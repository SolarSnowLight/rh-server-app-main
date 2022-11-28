package repository

import (
	//middlewareConstant "main-server/pkg/constant/middleware"
	"encoding/json"
	"errors"
	"fmt"
	actionConstant "main-server/pkg/constant/action"
	objectConstant "main-server/pkg/constant/object"
	roleConstant "main-server/pkg/constant/role"
	tableConstant "main-server/pkg/constant/table"
	projectModel "main-server/pkg/model/project"
	rbacModel "main-server/pkg/model/rbac"
	userModel "main-server/pkg/model/user"
	"main-server/pkg/model/worker"
	"sort"

	"strconv"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type ProjectPostgres struct {
	db       *sqlx.DB
	enforcer *casbin.Enforcer
	role     *RolePostgres
}

/* Function for create new struct of AdminPostgres */
func NewProjectPostgres(db *sqlx.DB, enforcer *casbin.Enforcer, role *RolePostgres) *ProjectPostgres {
	return &ProjectPostgres{
		db:       db,
		enforcer: enforcer,
		role:     role,
	}
}

func (r *ProjectPostgres) CreateProject(userId, domainId int, data projectModel.ProjectModel) (projectModel.ProjectModel, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return projectModel.ProjectModel{}, err
	}

	query := fmt.Sprintf("SELECT id FROM %s WHERE uuid=$1", tableConstant.COMPANIES_TABLE)
	var companyId int

	row := r.db.QueryRow(query, data.Uuid)
	if err := row.Scan(&companyId); err != nil {
		tx.Rollback()
		return projectModel.ProjectModel{}, err
	}

	query = fmt.Sprintf("INSERT INTO %s (uuid, data, created_at, updated_at, users_id, companies_id) values ($1, $2, $3, $4, $5, $6) RETURNING id", tableConstant.PROJECTS_TABLE)

	dataJson, err := json.Marshal(projectModel.ProjectDataModel{
		Title:       data.Title,
		Description: data.Description,
		Managers:    data.Managers,
	})

	if err != nil {
		tx.Rollback()
		return projectModel.ProjectModel{}, err
	}

	currentDate := time.Now()
	projectUuid := uuid.NewV4()

	var projectId int
	row = tx.QueryRow(query, projectUuid, dataJson, currentDate, currentDate, userId, companyId)
	if err := row.Scan(&projectId); err != nil {
		tx.Rollback()
		return projectModel.ProjectModel{}, err
	}

	if err != nil {
		tx.Rollback()
		return projectModel.ProjectModel{}, err
	}

	roleBuilderManager, err := r.role.GetRole("value", roleConstant.ROLE_BUILDER_MANAGER)
	if err != nil {
		tx.Rollback()
		return projectModel.ProjectModel{}, err
	}

	roleBuilderManagerJson, err := json.Marshal(worker.WorkerModel{
		Role: roleBuilderManager.Uuid,
	})

	if err != nil {
		tx.Rollback()
		return projectModel.ProjectModel{}, err
	}

	query = fmt.Sprintf(`
	INSERT INTO %s (uuid, data, created_at, updated_at, users_id, companies_id) 
	values ($1, $2, $3, $4, $5, $6) RETURNING id`,
		tableConstant.WORKERS_TABLE,
	)

	for _, element := range data.Managers {
		var user userModel.UserModel
		queryLocal := fmt.Sprintf("SELECT * FROM %s tl WHERE tl.email=$1", tableConstant.USERS_TABLE)
		err = r.db.Get(&user, queryLocal, element.Email)
		if err != nil {
			tx.Rollback()
			return projectModel.ProjectModel{}, err
		}

		queryFindWorker := fmt.Sprintf("SELECT * FROM %s tl WHERE tl.users_id = $1", tableConstant.WORKERS_TABLE)
		var workers []int
		err = r.db.Select(&workers, queryFindWorker, user.Id)

		var workerId int
		row = tx.QueryRow(query, uuid.NewV4().String(), roleBuilderManagerJson, currentDate, currentDate, user.Id, companyId)
		if err := row.Scan(&workerId); err != nil {
			tx.Rollback()
			return projectModel.ProjectModel{}, err
		}

		queryLocal = fmt.Sprintf(`
		INSERT INTO %s (workers_id, projects_id) 
		values ($1, $2)`,
			tableConstant.WORKERS_PROJECTS_TABLE,
		)

		_, err = tx.Exec(queryLocal, workerId, projectId)
		if err != nil {
			tx.Rollback()
			return projectModel.ProjectModel{}, err
		}
	}

	// Adding information about a resource
	var typesObjects rbacModel.TypesObjectsModel

	query = fmt.Sprintf("SELECT * FROM %s WHERE value=$1", tableConstant.TYPES_OBJECTS_TABLE)

	err = r.db.Get(&typesObjects, query, objectConstant.TYPE_PROJECT)
	if err != nil {
		tx.Rollback()
		return projectModel.ProjectModel{}, err
	}

	fmt.Println(typesObjects.Id)

	query = fmt.Sprintf(`INSERT INTO %s (value, types_objects_id) values ($1, $2)`, tableConstant.OBJECTS_TABLE)

	_, err = tx.Exec(query, projectUuid.String(), typesObjects.Id)
	if err != nil {
		tx.Rollback()
		return projectModel.ProjectModel{}, err
	}

	var userIdStr string = strconv.Itoa(userId)
	var domainIdStr string = strconv.Itoa(domainId)

	// Update current user policy for current article
	_, err = r.enforcer.AddPolicies([][]string{
		{userIdStr, domainIdStr, projectUuid.String(), actionConstant.DELETE},
		{userIdStr, domainIdStr, projectUuid.String(), actionConstant.MODIFY},
		{userIdStr, domainIdStr, projectUuid.String(), actionConstant.READ},
		{userIdStr, domainIdStr, projectUuid.String(), actionConstant.ADMINISTRATION},
	})

	if err != nil {
		tx.Rollback()
		return projectModel.ProjectModel{}, err
	}

	// Get ids for user managers
	query = fmt.Sprintf("SELECT id FROM %s WHERE email=$1 LIMIT 1", tableConstant.USERS_TABLE)

	for _, element := range data.Managers {
		// @(idea):
		// Можно улучшить данный фрагмент кода добавив сюда определение функции транзакции,
		// которая бы выполнялась в текущей транзакции

		var userManagerId int
		row := r.db.QueryRow(query, element.Email)
		if err := row.Scan(&userManagerId); err != nil {
			continue
		}

		userManagerIdStr := strconv.Itoa(userManagerId)

		_, err = r.enforcer.AddPolicies([][]string{
			{userManagerIdStr, domainIdStr, projectUuid.String(), actionConstant.READ},
			{userManagerIdStr, domainIdStr, projectUuid.String(), actionConstant.MANAGEMENT},
		})
	}

	// Save results all operation into a tables
	if err := tx.Commit(); err != nil {
		// Rejected
		r.enforcer.RemovePolicies([][]string{
			{userIdStr, domainIdStr, projectUuid.String(), actionConstant.DELETE},
			{userIdStr, domainIdStr, projectUuid.String(), actionConstant.MODIFY},
			{userIdStr, domainIdStr, projectUuid.String(), actionConstant.READ},
			{userIdStr, domainIdStr, projectUuid.String(), actionConstant.ADMINISTRATION},
		})

		for _, element := range data.Managers {
			var userManagerId int
			row := r.db.QueryRow(query, element.Email)
			if err := row.Scan(&userManagerId); err != nil {
				continue
			}

			userManagerIdStr := strconv.Itoa(userManagerId)

			_, err = r.enforcer.RemovePolicies([][]string{
				{userManagerIdStr, domainIdStr, projectUuid.String(), actionConstant.READ},
				{userManagerIdStr, domainIdStr, projectUuid.String(), actionConstant.MANAGEMENT},
			})
		}

		tx.Rollback()
		return projectModel.ProjectModel{}, err
	}

	return projectModel.ProjectModel{
		Logo:        data.Logo,
		Title:       data.Title,
		Description: data.Description,
		Managers:    data.Managers,
		Uuid:        projectUuid.String(),
	}, nil
}

/* Add logo project */
func (r *ProjectPostgres) ProjectUpdateImage(userId, domainId int, data projectModel.ProjectImageModel) (projectModel.ProjectImageModel, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return projectModel.ProjectImageModel{}, err
	}

	userIdStr := strconv.Itoa(userId)
	domainIdStr := strconv.Itoa(domainId)

	// Check access for user
	access, err := r.enforcer.Enforce(userIdStr, domainIdStr, data.Uuid, actionConstant.MODIFY)
	if err != nil {
		tx.Rollback()
		return projectModel.ProjectImageModel{}, err
	}

	if !access {
		tx.Rollback()
		return projectModel.ProjectImageModel{}, errors.New("Ошибка! Нет доступа!")
	}

	// Update logo for project
	query := fmt.Sprintf(`UPDATE %s tl SET data = jsonb_set(data, '{logo}', to_jsonb($1::text), true) WHERE tl.uuid = $2`, tableConstant.PROJECTS_TABLE)

	_, err = r.db.Exec(query, data.Filepath, data.Uuid)
	if err != nil {
		tx.Rollback()
		return projectModel.ProjectImageModel{}, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return projectModel.ProjectImageModel{}, err
	}

	return data, nil
}

/* Get information about one project */
func (r *ProjectPostgres) GetProject(userId, domainId int, data projectModel.ProjectUuidModel) (projectModel.ProjectDbModel, error) {
	var project projectModel.ProjectDbModel

	query := fmt.Sprintf("SELECT uuid, data, created_at FROM %s tl WHERE tl.uuid = $1 LIMIT 1", tableConstant.PROJECTS_TABLE)

	err := r.db.Get(&project, query, data.Uuid)
	if err != nil {
		return projectModel.ProjectDbModel{}, err
	}

	return project, nil
}

/* Получение информации обо всех проектах */
func (r *ProjectPostgres) GetProjects(userId, domainId int, data projectModel.ProjectCountModel) (projectModel.ProjectAnyCountModel, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE uuid=$1", tableConstant.COMPANIES_TABLE)
	var companyId int

	row := r.db.QueryRow(query, data.Uuid)
	if err := row.Scan(&companyId); err != nil {
		return projectModel.ProjectAnyCountModel{}, err
	}

	var projects []projectModel.ProjectDbModel
	sum := (data.Count + data.Limit)

	query = fmt.Sprintf("SELECT uuid, data, created_at FROM %s tl WHERE tl.companies_id = $1", tableConstant.PROJECTS_TABLE)
	err := r.db.Select(&projects, query, companyId)
	if err != nil {
		return projectModel.ProjectAnyCountModel{}, err
	}

	var projectsEx []projectModel.ProjectDbDataEx
	for _, element := range projects {
		var projectData projectModel.ProjectDataModel
		err = json.Unmarshal([]byte(element.Data), &projectData)

		if err != nil {
			return projectModel.ProjectAnyCountModel{}, err
		}

		projectsEx = append(projectsEx, projectModel.ProjectDbDataEx{
			Uuid:      element.Uuid,
			Data:      projectData,
			CreatedAt: element.CreatedAt,
		})
	}

	sort.SliceStable(projectsEx, func(i, j int) bool {
		return projectsEx[i].CreatedAt.After(projectsEx[j].CreatedAt)
	})

	if data.Count >= len(projectsEx) {
		return projectModel.ProjectAnyCountModel{}, nil
	}

	if sum >= len(projectsEx) {
		sum -= (sum - len(projectsEx))
	}

	return projectModel.ProjectAnyCountModel{
		Projects: projectsEx[data.Count:sum],
		Count:    (sum - data.Count),
	}, nil
}

/* Get information about any count managers for current project */
func (r *ProjectPostgres) GetProjectManagers(userId, domainId int, data projectModel.ProjectCountModel) (projectModel.ProjectAnyCountModel, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE uuid=$1", tableConstant.COMPANIES_TABLE)
	var companyId int

	row := r.db.QueryRow(query, data.Uuid)
	if err := row.Scan(&companyId); err != nil {
		return projectModel.ProjectAnyCountModel{}, err
	}

	var projects []projectModel.ProjectDbModel
	sum := (data.Count + data.Limit)

	query = fmt.Sprintf("SELECT uuid, data, created_at FROM %s tl WHERE tl.companies_id = $1 LIMIT $2", tableConstant.PROJECTS_TABLE)
	err := r.db.Select(&projects, query, companyId, sum)
	if err != nil {
		return projectModel.ProjectAnyCountModel{}, err
	}

	var projectsEx []projectModel.ProjectDbDataEx
	for _, element := range projects {
		var projectData projectModel.ProjectDataModel
		err = json.Unmarshal([]byte(element.Data), &projectData)

		if err != nil {
			return projectModel.ProjectAnyCountModel{}, err
		}

		projectsEx = append(projectsEx, projectModel.ProjectDbDataEx{
			Uuid:      element.Uuid,
			Data:      projectData,
			CreatedAt: element.CreatedAt,
		})
	}

	sort.SliceStable(projectsEx, func(i, j int) bool {
		return projectsEx[i].CreatedAt.After(projectsEx[j].CreatedAt)
	})

	if data.Count >= len(projectsEx) {
		return projectModel.ProjectAnyCountModel{}, nil
	}

	if sum >= len(projectsEx) {
		sum -= (sum - len(projectsEx))
	}

	return projectModel.ProjectAnyCountModel{
		Projects: projectsEx[data.Count:sum],
		Count:    (sum - data.Count),
	}, nil
}

func (r *ProjectPostgres) ProjectUpdate(user userModel.UserIdentityModel, data projectModel.ProjectUpdateModel) (projectModel.ProjectUpdateModel, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return projectModel.ProjectUpdateModel{}, err
	}

	userIdStr := strconv.Itoa(user.UserId)
	domainIdStr := strconv.Itoa(user.DomainId)

	// Check access for user
	access, err := r.enforcer.Enforce(userIdStr, domainIdStr, data.Uuid, actionConstant.MODIFY)
	if err != nil {
		tx.Rollback()
		return projectModel.ProjectUpdateModel{}, err
	}

	if !access {
		tx.Rollback()
		return projectModel.ProjectUpdateModel{}, errors.New("Ошибка! Нет доступа!")
	}

	var projectInfo []projectModel.ProjectDbModel
	query := fmt.Sprintf("SELECT uuid, data, created_at FROM %s WHERE uuid=$1 LIMIT 1", tableConstant.PROJECTS_TABLE)

	err = r.db.Select(&projectInfo, query, data.Uuid)
	if err != nil {
		return projectModel.ProjectUpdateModel{}, err
	}

	var projectData projectModel.ProjectDataModel
	err = json.Unmarshal([]byte(projectInfo[0].Data), &projectData)

	projectData.Description = data.Description
	projectData.Title = data.Title
	// projectData.Managers = data.Managers (need make)

	projectDataJson, err := json.Marshal(projectData)
	if err != nil {
		tx.Rollback()
		return projectModel.ProjectUpdateModel{}, err
	}

	query = fmt.Sprintf("UPDATE %s tl SET data=$1 WHERE tl.uuid=$2", tableConstant.PROJECTS_TABLE)

	_, err = tx.Exec(query, projectDataJson, data.Uuid)
	if err != nil {
		tx.Rollback()
		return projectModel.ProjectUpdateModel{}, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return projectModel.ProjectUpdateModel{}, err
	}

	return data, nil
}
