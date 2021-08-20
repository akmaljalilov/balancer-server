package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatusMaster(t *testing.T) {
	//StartCanal()
	//StatusMaster()
	//StartSync()
	//err := createUserForSlave2()
	//assert.NoError(t, err)
	//err = registerSlave()
	//assert.NoError(t, err)
	m, s, err := failoverCheck()
	assert.NoError(t, err)
	status, err := m.MasterStatus()
	assert.NoError(t, err)
	println(status)
	status, err = s.MasterStatus()
	println(status)
}
