package character_test

import (
	"context"
	dragonball "dragonball/business/core/dragon-ball"
	"dragonball/business/core/dragon-ball/character"
	"dragonball/foundation/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var log = logger.New(os.Stdout, logger.LevelDebug, "test", nil)

func TestCECLAP2(t *testing.T) {
	mock := newDragonBallMock()
	defer mock.Close()

	cfgPMC := dragonball.Config{
		BaseURL:          mock.URL,
		RetryCount:       2,
		RetryMaxWaitTime: time.Second * 1,
		Timeout:          time.Second * 30,
	}

	t.Parallel()

	t.Run("Run_Succed", RunSuccess(cfgPMC))

	t.Run("Run_Failed", RunFailed(cfgPMC))

	t.Run("Run_Error", RunError(cfgPMC))
}

func RunSuccess(cfg dragonball.Config) func(t *testing.T) {
	return func(t *testing.T) {
		name := "goku"
		dragonBallCore, err := dragonball.NewCore(log, cfg)
		assert.NoErrorf(t, err, "Should not return an error")

		_, err = character.Run(context.Background(), dragonBallCore, name)
		assert.NoError(t, err, "Given a correct input, it should not return an error")
	}
}

func RunFailed(cfg dragonball.Config) func(t *testing.T) {
	return func(t *testing.T) {
		name := "juan"
		dragonBallCore, err := dragonball.NewCore(log, cfg)
		assert.NoErrorf(t, err, "Should not return an error")

		dragonBallCore.Client.BaseURL += "/nodata"

		_, err = character.Run(context.Background(), dragonBallCore, name)
		assert.Errorf(t, err, "Given a correct input but no data is returned, it should return error")
	}
}

func RunError(cfg dragonball.Config) func(t *testing.T) {
	return func(t *testing.T) {
		dragonBallCore, err := dragonball.NewCore(log, cfg)
		assert.NoErrorf(t, err, "Should not return an error")

		dragonBallCore.Client.SetBaseURL("dragonball-api.com/api")

		_, err = character.Run(context.Background(), dragonBallCore, "")
		assert.Errorf(t, err, "Given a bad URL, it should return an error")
	}
}

// ======================================================================
//	VARS,CONST, AND HELPERS
// ======================================================================

func newDragonBallMock() *httptest.Server {
	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.Contains(r.URL.Path, "/nodata/characters"):
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(succeedWithNoData))
		case strings.Contains(r.URL.Path, "/characters"):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(success))
		}
	}

	return httptest.NewServer(http.HandlerFunc(h))
}

const (
	succeedWithNoData = ``
	success           = `[
		{
			"id": 1,
			"name": "Goku",
			"ki": "60.000.000",
			"maxKi": "90 Septillion",
			"race": "Saiyan",
			"gender": "Male",
			"description": "El protagonista de la serie, conocido por su gran poder y personalidad amigable. Originalmente enviado a la Tierra como un infante volador con la misión de conquistarla. Sin embargo, el caer por un barranco le proporcionó un brutal golpe que si bien casi lo mata, este alteró su memoria y anuló todos los instintos violentos de su especie, lo que lo hizo crecer con un corazón puro y bondadoso, pero conservando todos los poderes de su raza. No obstante, en la nueva continuidad de Dragon Ball se establece que él fue enviado por sus padres a la Tierra con el objetivo de sobrevivir a toda costa a la destrucción de su planeta por parte de Freeza. Más tarde, Kakarot, ahora conocido como Son Goku, se convertiría en el príncipe consorte del monte Fry-pan y líder de los Guerreros Z, así como el mayor defensor de la Tierra y del Universo 7, logrando mantenerlos a salvo de la destrucción en innumerables ocasiones, a pesar de no considerarse a sí mismo como un héroe o salvador.",
			"image": "https://dragonball-api.com/characters/goku_normal.webp",
			"affiliation": "Z Fighter",
			"deletedAt": null
		}
	]`
)
