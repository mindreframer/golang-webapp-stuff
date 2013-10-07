package mgo

import (
	"fmt"
	. "github.com/ku-ovdp/api/entities"
	. "github.com/ku-ovdp/api/repository"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func (m *mgoBackend) NewProjectRepository(repositories RepositoryGroup) ProjectRepository {
	if m.db == nil {
		panic("Not initialized! Missing call to Init()?")
	}
	return &projectRepository{m, m.db.C("projects"), repositories}
}

type projectRepository struct {
	backend *mgoBackend
	c       *mgo.Collection
	group   RepositoryGroup
}

func (pr projectRepository) Get(id int) (Project, error) {
	result := Project{}
	q := pr.c.Find(bson.M{"id": id})
	err := q.One(&result)
	if err == mgo.ErrNotFound {
		return result, NewErrNotFound(Project{}, id)
	}
	return result, err
}

func (pr projectRepository) Put(project Project) (Project, error) {
	if project.Id == 0 {
		project.Id = pr.backend.nextId("project")
	}
	_, err := pr.c.Upsert(bson.M{"id": project.Id}, project)
	if err != nil {
		return Project{}, err
	}
	return pr.Get(project.Id)
}

func (pr projectRepository) Remove(id int) error {
	return fmt.Errorf("Not implemented")
}

func (pr projectRepository) Scan(from, to int) ([]Project, error) {
	query := bson.M{"id": bson.M{"$gte": from}}
	if to > 0 {
		query["id"].(bson.M)["$lte"] = to
	}

	q := pr.c.Find(query)

	results := []Project{}
	err := q.All(&results)

	return results, err
}

func (pr projectRepository) Group() RepositoryGroup {
	return pr.group
}
