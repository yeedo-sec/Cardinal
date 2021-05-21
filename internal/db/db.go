// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by Apache-2.0
// license that can be found in the LICENSE file.

package db

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/vidar-team/Cardinal/internal/dbold"
)

var ErrBadCharset = errors.New("bad charset")

// Init initializes the database.
func Init(username, password, host, name string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Local&charset=utf8mb4,utf8", username, password, host, name)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return errors.Wrap(err, "open connection")
	}

	// Migrate databases.
	if db.AutoMigrate(
		&Challenge{},
	) != nil {
		return errors.Wrap(err, "auto migrate")
	}

	// Test database charset, we should support Chinese input.
	if dbold.MySQL.Exec("SELECT * FROM `logs` WHERE `Content` = '中文测试';").Error != nil {
		return ErrBadCharset
	}

	Challenges = NewChallengesStore(db)

	return nil
}