package controllers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"vms/internal/controllers"
	"vms/internal/models"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestDatabase() error {
	// Set up the connection to the test database
	err := reposity.Connect("localhost", "5432", "test_db", "disable", "postgres", "123", "public")
	if err != nil {
		return err
	}

	// Migrate the schema (create the necessary tables)
	err = reposity.Migrate(&models.AIWaring{}) // Add the relevant model(s)
	if err != nil {
		return err
	}

	return nil
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ExecWithPaging(sort string, limit, page int) ([]models.DTO_AI_Event, int, error) {
	args := m.Called(sort, limit, page)
	return args.Get(0).([]models.DTO_AI_Event), args.Int(1), args.Error(2)
}

func insertSampleEntry() (*models.DTO_AI_Event, error) {
	// Create a sample entry in the database
	sampleEvent := models.DTO_AI_Event{
		ID:              uuid.New(),
		CamName:         "Camera 1",
		EventTypeString: "Motion",
		// Add other fields if necessary
	}

	createdEvent, err := reposity.CreateItemFromDTO[models.DTO_AI_Event, models.AIWaring](sampleEvent)
	if err != nil {
		return nil, err
	}

	return &createdEvent, nil
}

func TestSearchAIEvent(t *testing.T) {
	err := setupTestDatabase()
	if err != nil {
		log.Fatal("Failed to set up test database:", err)
	}
	t.Run("Success with data", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Simulate a valid request
		c.Request = httptest.NewRequest("GET", "/aievents/searching?keyword=Motion&limit=25&page=1", nil)

		// Call the handler
		controllers.SearchAIEvent(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Camera 1") // Assuming "Camera 1" is expected in the output
	})

	t.Run("No data scenario", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Simulate a valid request but expect no data
		c.Request = httptest.NewRequest("GET", "/aievents/searching?keyword=nonexistent", nil)

		// Call the handler
		controllers.SearchAIEvent(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"data":[]`) // Expecting empty data
	})
}

func TestGetAIEvent(t *testing.T) {
	err := setupTestDatabase()
	if err != nil {
		log.Fatal("Failed to set up test database:", err)
	}
	t.Run("Success with filters", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/aievents?camName=Cam1&limit=25&page=1", nil)

		controllers.GetAIEvent(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid time range", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/aievents?startTime=invalid&endTime=invalid", nil)
		controllers.GetAIEvent(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "invalid input syntax")
	})

	t.Run("No data scenario", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/aievents?camName=nonexistent&limit=25&page=1", nil)

		// Call the handler
		controllers.GetAIEvent(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"data":[]`)
	})
}

func TestColumnChartDay(t *testing.T) {
	err := setupTestDatabase()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	createRequestContext := func(queryParams string) (*gin.Context, *httptest.ResponseRecorder) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req := httptest.NewRequest("GET", "/aievents/columnchartday?"+queryParams, nil)
		c.Request = req
		return c, w
	}

	t.Run("Valid data", func(t *testing.T) {
		c, w := createRequestContext("startTime=1627891200&endTime=1627894800")

		controllers.ColumnChartDay(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
		assert.Contains(t, w.Body.String(), "Success")
		assert.Contains(t, w.Body.String(), "Success")
	})

	t.Run("No matching data", func(t *testing.T) {
		c, w := createRequestContext("startTime=0&endTime=0")
		controllers.ColumnChartDay(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"data":[]`)
	})

	t.Run("Invalid time format", func(t *testing.T) {
		c, w := createRequestContext("startTime=invalid&endTime=1627894800")
		controllers.ColumnChartDay(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid startTime")
	})

	t.Run("Partial matching data", func(t *testing.T) {
		c, w := createRequestContext("startTime=1627891200&endTime=1627894800")
		controllers.ColumnChartDay(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
		assert.Contains(t, w.Body.String(), "Success")
		assert.Contains(t, w.Body.String(), "Success")
	})
}

func TestUpdateAIWarning(t *testing.T) {
	err := setupTestDatabase()
	if err != nil {
		log.Fatal("Failed to set up test database:", err)
	}

	t.Run("Successful update", func(t *testing.T) {
		sampleEvent, err := insertSampleEntry()
		if err != nil {
			t.Fatalf("Failed to insert sample entry: %v", err)
		}

		r := gin.Default()
		r.PUT("/aievents/update/:id", controllers.UpdateAIWarning)

		w := httptest.NewRecorder()
		jsonBody := `{"typeOfAIEvent": "New Event"}`
		req := httptest.NewRequest("PUT", "/aievents/update/"+sampleEvent.ID.String(), strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
	})

	t.Run("Invalid JSON in update", func(t *testing.T) {
		r := gin.Default()
		r.PUT("/aievents/update/:id", controllers.UpdateAIWarning)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/aievents/update/valid-id", strings.NewReader(`invalid-json`))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid character")
	})
}

func TestAIEventRoutineHandler(t *testing.T) {
	err := setupTestDatabase()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	t.Run("Success with valid filters", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/aievents/routine?keyword=test&eventType=AI_EVENT_BLACKLIST_FACE_RECOGNITION&startTime=1627891200&endTime=1627894800", nil)
		c.Request = req

		controllers.AIEventRoutineHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Success")
		assert.Contains(t, w.Body.String(), `"data"`)
	})

	t.Run("No data with valid filters", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/aievents/routine?keyword=emptykeyword", nil)
		c.Request = req

		controllers.AIEventRoutineHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"data":[]`)
	})

	t.Run("Invalid time format in query params", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/aievents/routine?startTime=invalid-time", nil)
		c.Request = req

		controllers.AIEventRoutineHandler(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "invalid input syntax")
	})

	// t.Run("Database error", func(t *testing.T) {
	// 	w := httptest.NewRecorder()
	// 	c, _ := gin.CreateTestContext(w)

	// 	req := httptest.NewRequest("GET", "/aievents/routine?keyword=test", nil)
	// 	c.Request = req

	// 	controllers.AIEventRoutineHandler(c)
	// 	assert.Equal(t, http.StatusInternalServerError, w.Code)
	// 	assert.Contains(t, w.Body.String(), "database error")
	// })

	t.Run("Empty request with default values", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/aievents/routine", nil)
		c.Request = req

		controllers.AIEventRoutineHandler(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"data"`)
	})
}

func TestReadCabinEventAI(t *testing.T) {
	err := setupTestDatabase()
	if err != nil {
		log.Fatal("Failed to set up test database:", err)
	}

	createRequestContext := func(pathParam string) (*gin.Context, *httptest.ResponseRecorder) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = append(c.Params, gin.Param{Key: "id", Value: pathParam})
		return c, w
	}

	t.Run("Success case - valid ID", func(t *testing.T) {
		existingEntry, err := reposity.CreateItemFromDTO[models.DTO_AI_Event, models.AIWaring](models.DTO_AI_Event{
			EventType: "AI_EVENT_MOTION",
			CamName:   "Test Camera",
			CameraId:  "camera-1",
			Location:  "Test Location",
			Timestamp: time.Now().Unix(),
		})

		if err != nil {
			log.Fatalf("Failed to create valid cabin event for test: %v", err)
		}

		validID := existingEntry.ID

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{{Key: "id", Value: validID.String()}}

		controllers.ReadCabinEventAI(c)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"Data"`)
	})

	t.Run("Not found case - invalid ID", func(t *testing.T) {
		c, w := createRequestContext("invalid-id")

		controllers.ReadCabinEventAI(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "")
	})

	t.Run("Empty ID", func(t *testing.T) {
		c, w := createRequestContext("")

		controllers.ReadCabinEventAI(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "")
	})
}
