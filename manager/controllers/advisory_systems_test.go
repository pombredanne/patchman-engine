package controllers

import (
	"app/base/core"
	"app/base/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdvisorySystemsDefault(t *testing.T) { //nolint:dupl
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/RH-1", nil)
	core.InitRouterWithPath(AdvisorySystemsListHandler, "/:advisory_id").ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var output AdvisorySystemsResponse
	ParseReponseBody(t, w.Body.Bytes(), &output)
	assert.Equal(t, 8, len(output.Data))
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", output.Data[0].ID)
	assert.Equal(t, "system", output.Data[0].Type)
	assert.Equal(t, "2020-09-22 16:00:00 +0000 UTC", output.Data[0].Attributes.LastUpload.String())
	assert.Equal(t, 2, output.Data[0].Attributes.RhsaCount)
	assert.Equal(t, 3, output.Data[0].Attributes.RhbaCount)
	assert.Equal(t, 3, output.Data[0].Attributes.RheaCount)
}

func TestAdvisorySystemsNotFound(t *testing.T) { //nolint:dupl
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nonexistant/systems", nil)
	core.InitRouterWithPath(AdvisorySystemsListHandler, "/:advisory_id").ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAdvisorySystemsOffsetLimit(t *testing.T) { //nolint:dupl
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/RH-1?offset=5&limit=3", nil)
	core.InitRouterWithPath(AdvisorySystemsListHandler, "/:advisory_id").ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var output AdvisorySystemsResponse
	ParseReponseBody(t, w.Body.Bytes(), &output)
	assert.Equal(t, 3, len(output.Data))
	assert.Equal(t, "00000000-0000-0000-0000-000000000006", output.Data[0].ID)
	assert.Equal(t, "system", output.Data[0].Type)
	assert.Equal(t, "2018-08-26 16:00:00 +0000 UTC", output.Data[0].Attributes.LastUpload.String())
	assert.Equal(t, 0, output.Data[0].Attributes.RhsaCount)
	assert.Equal(t, 1, output.Data[0].Attributes.RheaCount)
	assert.Equal(t, 0, output.Data[0].Attributes.RhbaCount)
}

func TestAdvisorySystemsOffsetOverflow(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/RH-1?offset=100&limit=3", nil)
	core.InitRouterWithPath(AdvisorySystemsListHandler, "/:advisory_id").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp utils.ErrorResponse
	ParseReponseBody(t, w.Body.Bytes(), &errResp)
	assert.Equal(t, InvalidOffsetMsg, errResp.Error)
}

func TestAdvisorySystemsSorts(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	for sort := range SystemsFields {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/RH-1?sort=%v", sort), nil)
		core.InitRouterWithPath(AdvisorySystemsListHandler, "/:advisory_id").ServeHTTP(w, req)

		var output AdvisorySystemsResponse
		ParseReponseBody(t, w.Body.Bytes(), &output)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, 1, len(output.Meta.Sort))
		assert.Equal(t, output.Meta.Sort[0], sort)
	}
}

func TestAdvisorySystemsWrongSort(t *testing.T) { //nolint:dupl
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/RH-1?sort=unknown_key", nil)
	core.InitRouterWithPath(AdvisorySystemsListHandler, "/:advisory_id").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAdvisorySystemsTags(t *testing.T) { //nolint:dupl
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/RH-1?tags=ns1/k1=val1", nil)
	core.InitRouterWithPath(AdvisorySystemsListHandler, "/:advisory_id").ServeHTTP(w, req)

	var output AdvisorySystemsResponse
	ParseReponseBody(t, w.Body.Bytes(), &output)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 7, len(output.Data))
}

func TestAdvisorySystemsTagsMultiple(t *testing.T) { //nolint:dupl
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/RH-1?tags=ns1/k3=val4&tags=ns1/k1=val1", nil)
	core.InitRouterWithPath(AdvisorySystemsListHandler, "/:advisory_id").ServeHTTP(w, req)

	var output AdvisorySystemsResponse
	ParseReponseBody(t, w.Body.Bytes(), &output)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 1, len(output.Data))
	assert.Equal(t, "00000000-0000-0000-0000-000000000003", output.Data[0].ID)
}

func TestAdvisorySystemsTagsInvalid(t *testing.T) {
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/RH-1?tags=ns1/k3=val4&tags=invalidTag", nil)
	core.InitRouterWithPath(AdvisorySystemsListHandler, "/:advisory_id").ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp utils.ErrorResponse
	ParseReponseBody(t, w.Body.Bytes(), &errResp)
	assert.Equal(t, fmt.Sprintf(InvalidTagMsg, "invalidTag"), errResp.Error)
}

func TestAdvisorySystemsTagsUnknown(t *testing.T) { //nolint:dupl
	utils.SkipWithoutDB(t)
	core.SetupTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/RH-1?tags=ns1/k3=val4&tags=ns1/k1=unk", nil)
	core.InitRouterWithPath(AdvisorySystemsListHandler, "/:advisory_id").ServeHTTP(w, req)

	var output AdvisorySystemsResponse
	ParseReponseBody(t, w.Body.Bytes(), &output)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 0, len(output.Data))
}

func TestAdvisorySystemsWrongOffset(t *testing.T) {
	doTestWrongOffset(t, "/:advisory_id", "/RH-1?offset=1000", AdvisorySystemsListHandler)
}
