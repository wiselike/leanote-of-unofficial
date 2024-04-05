package tests

import (
	"github.com/wiselike/leanote-of-unofficial/app/db"
	"testing"
	//	. "github.com/wiselike/leanote-of-unofficial/app/lea"
	//	"github.com/wiselike/leanote-of-unofficial/app/service"
	//	"gopkg.in/mgo.v2"
	//	"fmt"
)

func TestDBConnect(t *testing.T) {
	db.Init("mongodb://localhost:27017/leanote", "leanote")
}
