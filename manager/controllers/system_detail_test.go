package controllers

import (
	"app/base/core"
	"app/base/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSystemDetailDefault1(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/00000000-0000-0000-0000-000000000001", nil)
	core.InitRouterWithPath(SystemDetailHandler, "/:inventory_id").ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var output SystemDetailResponse
	ParseReponseBody(t, w.Body.Bytes(), &output)
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", output.Data.ID)
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", output.Data.Attributes.DisplayName)
	assert.Equal(t, "system", output.Data.Type)
	assert.Equal(t, "2018-09-22 16:00:00 +0000 UTC", output.Data.Attributes.LastEvaluation.String())
	assert.Equal(t, "2020-09-22 16:00:00 +0000 UTC", output.Data.Attributes.LastUpload.String())
	assert.False(t, output.Data.Attributes.Stale)
	assert.Equal(t, 2, output.Data.Attributes.RhsaCount)
	assert.Equal(t, 3, output.Data.Attributes.RhbaCount)
	assert.Equal(t, 3, output.Data.Attributes.RheaCount)
}

func TestSystemDetailDefault2(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	// get system with some installable/updatable packages
	req, _ := http.NewRequest("GET", "/00000000-0000-0000-0000-000000000012", nil)
	core.InitRouterWithAccount(SystemDetailHandler, "/:inventory_id", 3).ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var output SystemDetailResponse
	ParseReponseBody(t, w.Body.Bytes(), &output)
	assert.Equal(t, 2, output.Data.Attributes.PackagesInstalled)
	assert.Equal(t, 2, output.Data.Attributes.PackagesUpdatable)
}

func TestSystemDetailNoIdProvided(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	core.InitRouter(SystemDetailHandler).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp utils.ErrorResponse
	ParseReponseBody(t, w.Body.Bytes(), &errResp)
	assert.Equal(t, "inventory_id param not found", errResp.Error)
}

func TestSystemDetailNotFound(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	core.InitRouterWithPath(SystemDetailHandler, "/:inventory_id").ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	var errResp utils.ErrorResponse
	ParseReponseBody(t, w.Body.Bytes(), &errResp)
	assert.Equal(t, "inventory not found", errResp.Error)
}
