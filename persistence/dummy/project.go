// Package dummy implements dummy storage for API entities
package dummy

import (
	. "github.com/ku-ovdp/api/entities"
	. "github.com/ku-ovdp/api/repository"
	"time"
)

func (d dummyBackend) NewProjectRepository(repositories RepositoryGroup) ProjectRepository {
	return projectRepository{dummyProjectData, repositories}
}

type projectRepo map[int]Project

type projectRepository struct {
	repo  projectRepo
	group RepositoryGroup
}

var dummyProjectData = projectRepo{
	1: {Id: 1, Name: "Project Parkinson's",
		Slug:                   "parkinsons",
		HighlevelDescription:   "Project Parkinson's ... (high level)",
		DetailedDescription:    "Project Parkinson's ... (detailed)",
		PrivacyPolicyURL:       "http://openvoicedata.org/privacy.php",
		MinimumNumberOfSamples: 2,
		MaximumNumberOfSamples: 3,
		GeneralInstructions:    "(general instructions)",
		SampleInstructions: []SampleInstruction{
			{Duration: 10, Instruction: "Produce an 'Ah' sound at a comfortable level."},
			{Duration: 10, Instruction: "Produce an 'Ah' sound with twice the previous effort."},
			{Duration: 0, Instruction: "Produce normal conversational speaking"},
		},
		FormFields: []FormField{
			{Label: "Age", Slug: "age", Type: "int", Required: true, Description: "Your Age"},
			{Label: "Gender", Slug: "gender", Type: "string", Required: true, Description: "Your Gender",
				Meta: `{"options": ["Male", "Female", "Undisclosed"]}`},
			{Label: "Parkinson's Diagnosis", Slug: "parkinsons", Type: "bool", Required: true, Description: "Have you been diagnosed with Parkinson's?"},
		},
		Created: time.Now().Add(time.Hour * -24 * 14)},

	2: {Id: 2, Name: "Disphonia Foobar",
		Slug:    "foobar",
		Created: time.Now().Add(time.Hour * -24 * 10)},
}

func (pr projectRepository) Get(id int) (Project, error) {
	if p, ok := pr.repo[id]; ok {
		return p, nil
	} else {
		return Project{}, NewErrNotFound(Project{}, id)
	}
}

func (pr projectRepository) Put(project Project) (Project, error) {
	pr.repo[project.Id] = project
	return project, nil
}

func (pr projectRepository) Remove(id int) error {
	delete(pr.repo, id)
	return nil
}

func (pr projectRepository) Scan(from, to int) ([]Project, error) {
	results := []Project{}
	for id, value := range pr.repo {
		if id < from {
			continue
		}
		if id > to && to != 0 {
			continue
		}
		results = append(results, value)
	}
	return results, nil
}

func (pr projectRepository) Group() RepositoryGroup {
	return pr.group
}
