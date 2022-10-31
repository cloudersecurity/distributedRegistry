package grades

// 给学生们赋值
func init() {
	students = []Student{
		{
			ID:        1,
			FirstName: "Nick",
			LastName:  "Carter",
			Grades: []Grade{
				{
					Title: "Quiz1",
					Type:  GradeQuiz,
					Score: 85,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 94,
				},
				{
					Title: "Quiz2",
					Type:  GradeQuiz,
					Score: 82,
				},
			},
		},
		{
			ID:        2,
			FirstName: "Rob",
			LastName:  "Cathy",
			Grades: []Grade{
				{
					Title: "Quiz1",
					Type:  GradeQuiz,
					Score: 100,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 100,
				},
				{
					Title: "Quiz2",
					Type:  GradeQuiz,
					Score: 81,
				},
			},
		},
		{
			ID:        3,
			FirstName: "Emma",
			LastName:  "Stone",
			Grades: []Grade{
				{
					Title: "Quiz1",
					Type:  GradeQuiz,
					Score: 67,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 0,
				},
				{
					Title: "Quiz2",
					Type:  GradeQuiz,
					Score: 75,
				},
			},
		},
		{
			ID:        4,
			FirstName: "Rachel",
			LastName:  "Mike",
			Grades: []Grade{
				{
					Title: "Quiz1",
					Type:  GradeQuiz,
					Score: 98,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 99,
				},
				{
					Title: "Quiz2",
					Type:  GradeQuiz,
					Score: 94,
				},
			},
		},
	}
}
