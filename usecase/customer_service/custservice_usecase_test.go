package usecase

import (
	"calmind/model"
	"testing"
)

// Define a simple in-memory implementation of the repository for testing.
type InMemoryCustServiceRepo struct {
	Data []model.CustService
}

func (repo *InMemoryCustServiceRepo) SaveCustService(custService *model.CustService) error {
	repo.Data = append(repo.Data, *custService)
	return nil
}

func (repo *InMemoryCustServiceRepo) AnswerMessage(id int, answer string) error {
	for i, v := range repo.Data {
		if v.ID == id {
			repo.Data[i].Answer = answer
			repo.Data[i].IsAnswered = true
			return nil
		}
	}
	return nil
}

// Unit tests for SaveCustService and AnswerMessage methods.
func TestSaveCustService(t *testing.T) {
	// Arrange
	repo := &InMemoryCustServiceRepo{}
	usecase := NewCustServiceUsecase(repo)

	userID := 1
	message := "Test message"

	// Act
	err := usecase.SaveCustService(userID, message)

	// Assert
	if err != nil {
		t.Errorf("SaveCustService() failed, expected no error, got %v", err)
	}
	if len(repo.Data) != 1 {
		t.Errorf("SaveCustService() failed, expected 1 item in the repository, got %d", len(repo.Data))
	}
	if repo.Data[0].Message != message {
		t.Errorf("SaveCustService() failed, expected message %s, got %s", message, repo.Data[0].Message)
	}
	if repo.Data[0].UserID != userID {
		t.Errorf("SaveCustService() failed, expected userID %d, got %d", userID, repo.Data[0].UserID)
	}
	if repo.Data[0].IsAnswered {
		t.Errorf("SaveCustService() failed, expected IsAnswered to be false, got true")
	}
}

func TestAnswerMessage(t *testing.T) {
	// Arrange
	repo := &InMemoryCustServiceRepo{
		Data: []model.CustService{
			{ID: 1, Message: "Test message", IsAnswered: false},
		},
	}
	usecase := NewCustServiceUsecase(repo)

	id := 1
	answer := "Test answer"

	// Act
	err := usecase.AnswerMessage(id, answer)

	// Assert
	if err != nil {
		t.Errorf("AnswerMessage() failed, expected no error, got %v", err)
	}
	if repo.Data[0].Answer != answer {
		t.Errorf("AnswerMessage() failed, expected answer %s, got %s", answer, repo.Data[0].Answer)
	}
	if !repo.Data[0].IsAnswered {
		t.Errorf("AnswerMessage() failed, expected IsAnswered to be true, got false")
	}
}

func TestAnswerMessageWithInvalidID(t *testing.T) {
	// Arrange
	repo := &InMemoryCustServiceRepo{
		Data: []model.CustService{
			{ID: 1, Message: "Test message", IsAnswered: false},
		},
	}
	usecase := NewCustServiceUsecase(repo)

	invalidID := 2
	answer := "Test answer"

	// Act
	err := usecase.AnswerMessage(invalidID, answer)

	// Assert
	if err != nil {
		t.Errorf("AnswerMessage() failed, expected no error when ID is invalid, got %v", err)
	}
	if repo.Data[0].IsAnswered {
		t.Errorf("AnswerMessage() failed, expected IsAnswered to remain false, got true")
	}
	if repo.Data[0].Answer != "" {
		t.Errorf("AnswerMessage() failed, expected Answer to remain empty, got %s", repo.Data[0].Answer)
	}
}
