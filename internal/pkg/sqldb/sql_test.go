package sqldb

import (
	"context"
	"testing"
	"time"

	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/require"
)

// func TestFailedConnect(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
// 	defer cancel()

// 	// creates db config
// 	cfg := DBConfig{
// 		Driver:                "postgres",
// 		MasterDSN:             "localhost:5432",
// 		FollowerDSN:           "localhost:5432",
// 		ConnectionMaxLifetime: 1 * time.Minute,
// 		MaxIdleConnections:    2,
// 		MaxOpenConnections:    20,
// 		Retry:                 2,
// 	}

// 	// connect
// 	db, err := Connect(ctx, cfg)
// 	require.Error(t, err)
// 	require.Nil(t, db)
// }

func TestConnect(t *testing.T) {
	testCases := []struct {
		name string
		cfg  DBConfig
	}{
		{
			name: "ping check",
			cfg: DBConfig{
				Driver:                "ramsql",
				MasterDSN:             "dsn",
				FollowerDSN:           "dsn",
				ConnectionMaxLifetime: 1 * time.Minute,
				MaxIdleConnections:    2,
				MaxOpenConnections:    20,
				NoPingOnOpen:          false,
			},
		},
		{
			name: "no ping check",
			cfg: DBConfig{
				Driver:                "ramsql",
				MasterDSN:             "dsn",
				ConnectionMaxLifetime: 1 * time.Minute,
				MaxIdleConnections:    2,
				MaxOpenConnections:    20,
				NoPingOnOpen:          true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			// connect
			db, err := Connect(ctx, tc.cfg)
			require.NoError(t, err)
			require.NotNil(t, db)

			err = db.Ping()
			require.NoError(t, err)
		})
	}
}

// func TestConnectServerDown(t *testing.T) {
// 	testCases := []struct {
// 		name        string
// 		noPingCheck bool
// 		gotErr      bool
// 	}{
// 		{
// 			name:        "ping check",
// 			noPingCheck: false,
// 			gotErr:      true,
// 		},
// 		{
// 			name:        "no ping check",
// 			noPingCheck: true,
// 			gotErr:      false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctx := context.Background()

// 			// creates db config
// 			cfg := DBConfig{
// 				Driver:                "postgres",
// 				MasterDSN:             "postgres://toped:toped@localhost:1111?sslmode=disable",
// 				FollowerDSN:           "postgres://toped:toped@localhost:1111?sslmode=disable",
// 				ConnectionMaxLifetime: 1 * time.Minute,
// 				MaxIdleConnections:    2,
// 				MaxOpenConnections:    20,
// 				NoPingOnOpen:          tc.noPingCheck,
// 			}

// 			// connect
// 			_, err := Connect(ctx, cfg)
// 			if tc.gotErr {
// 				require.Error(t, err)
// 			} else {
// 				require.NoError(t, err)
// 			}
// 		})
// 	}
// }

//
func TestSelect(t *testing.T) {
	ctx := context.Background()

	// creates db config
	cfg := DBConfig{
		Driver:                "ramsql",
		MasterDSN:             "dsn",
		FollowerDSN:           "dsn",
		ConnectionMaxLifetime: 1 * time.Minute,
		MaxIdleConnections:    2,
		MaxOpenConnections:    20,
	}

	// connect
	db, err := Connect(ctx, cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	err = db.Ping()
	require.NoError(t, err)

	// create table
	_, err = db.Exec(`CREATE TABLE address (id BIGSERIAL PRIMARY KEY, street TEXT, street_number INT);`)
	require.NoError(t, err)

	const (
		street       = "hugo"
		streetNumber = 32
	)
	// insert
	q := db.Rebind("INSERT INTO address (street, street_number) VALUES ($1, $2)")
	_, err = db.Exec(q, street, streetNumber)
	require.NoError(t, err)

	// get
	var data struct {
		Street       string `db:"street"`
		StreetNumber int    `db:"street_number"`
	}

	err = db.Get(&data, `select street, street_number from address where street='hugo'`)
	require.NoError(t, err)
	require.Equal(t, street, data.Street)
	require.Equal(t, streetNumber, data.StreetNumber)
}

func TestObject(t *testing.T) {
	ctx := context.Background()

	// creates db config
	cfg := DBConfig{
		Driver:      "ramsql",
		MasterDSN:   "dsn",
		FollowerDSN: "dsn",
	}

	// connect
	db, err := Connect(ctx, cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	// test the master & follower getter
	master := db.GetMaster()
	require.Equal(t, db.master, master)

	follower := db.GetFollower()
	require.Equal(t, db.follower, follower)

	// test the statement
	writeStmt, err := db.PrepareWrite(ctx, "")
	require.NoError(t, err)
	require.Implements(t, (*WriteStatement)(nil), writeStmt)

	readStmt, err := db.PrepareRead(ctx, "")
	require.NoError(t, err)
	require.Implements(t, (*ReadStatement)(nil), readStmt)

}

func Test_getNoPassDSN(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		want string
	}{
		{
			name: "form1 mid",
			dsn:  "host=host1 user=user1 password='password1' dbname=db1 sslmode=disable",
			want: "host=host1 user=user1 dbname=db1 sslmode=disable",
		},
		{
			name: "form1 front",
			dsn:  "password='password1' host=host1 user=user1 dbname=db1 sslmode=disable",
			want: "host=host1 user=user1 dbname=db1 sslmode=disable",
		},
		{
			name: "form1 back",
			dsn:  "host=host1 user=user1 dbname=db1 sslmode=disable password='password1'",
			want: "host=host1 user=user1 dbname=db1 sslmode=disable",
		},
		{
			name: "form2",
			dsn:  "postgres://user2:password2@host2/db2?sslmode=disable",
			want: "postgres://user2@host2/db2?sslmode=disable",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNoPassDSN(tt.dsn); got != tt.want {
				t.Errorf("getNoPassDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Rebind1(t *testing.T) {
	type usecase struct {
		input    string
		expected string
	}
	tests := []usecase{
		{
			input:    "SELECT * FROM table_1 WHERE status = 1",
			expected: "SELECT * FROM table_1 WHERE status = 1",
		},
		{
			input:    "SELECT * FROM table_1 WHERE status = ?",
			expected: "SELECT * FROM table_1 WHERE status = $1",
		},
		{
			input:    "SELECT * FROM table_1 WHERE status = ? and id = ?",
			expected: "SELECT * FROM table_1 WHERE status = $1 and id = $2",
		},
		{
			input:    "SELECT * FROM table_1 WHERE status = ? and id = ? and amount = ?",
			expected: "SELECT * FROM table_1 WHERE status = $1 and id = $2 and amount = $3",
		},
	}

	db := &DB{driver: "postgres"}
	for _, v := range tests {
		result := db.Rebind(v.input)
		if result != v.expected {
			t.Errorf("Rebind() = %s, want %s", result, v.expected)
		}
	}
}

func Test_Rebind2(t *testing.T) {
	type usecase struct {
		input    string
		expected string
	}
	tests := []usecase{
		{
			input:    "SELECT * FROM table_1 WHERE status = 1",
			expected: "SELECT * FROM table_1 WHERE status = 1",
		},
		{
			input:    "SELECT * FROM table_1 WHERE status = ?",
			expected: "SELECT * FROM table_1 WHERE status = ?",
		},
		{
			input:    "SELECT * FROM table_1 WHERE status = ? and id = ?",
			expected: "SELECT * FROM table_1 WHERE status = ? and id = ?",
		},
		{
			input:    "SELECT * FROM table_1 WHERE status = ? and id = ? and amount = ?",
			expected: "SELECT * FROM table_1 WHERE status = ? and id = ? and amount = ?",
		},
	}

	db := &DB{driver: "mysql"}
	for _, v := range tests {
		result := db.Rebind(v.input)
		if result != v.expected {
			t.Errorf("Rebind() = %s, want %s", result, v.expected)
		}
	}
}

func Test_BindNamed1(t *testing.T) {
	type usecase struct {
		input    string
		params   map[string]interface{}
		expected string
	}
	tests := []usecase{
		{
			input: "SELECT * FROM table_1 WHERE status = :param1",
			params: map[string]interface{}{
				"param1": 1,
			},
			expected: "SELECT * FROM table_1 WHERE status = $1",
		},
		{
			input: "SELECT * FROM table_1 WHERE status = :param1 and id = :param2",
			params: map[string]interface{}{
				"param1": 1,
				"param2": 2,
			},
			expected: "SELECT * FROM table_1 WHERE status = $1 and id = $2",
		},
		{
			input: "SELECT * FROM table_1 WHERE status = :param1 and id = :param2 and amount = :param3",
			params: map[string]interface{}{
				"param1": 1,
				"param2": 2,
				"param3": 3,
			},
			expected: "SELECT * FROM table_1 WHERE status = $1 and id = $2 and amount = $3",
		},
	}

	db := &DB{driver: "postgres"}
	for _, v := range tests {
		result, _, err := db.BindNamed(v.input, v.params)
		if err != nil {
			t.Errorf("BindNamed() = %s, want not error", err)
		}
		if result != v.expected {
			t.Errorf("BindNamed() = %s, want %s", result, v.expected)
		}
	}
}

func Test_BindNamed2(t *testing.T) {
	type usecase struct {
		input    string
		params   map[string]interface{}
		expected string
	}
	tests := []usecase{
		{
			input: "SELECT * FROM table_1 WHERE status = :param1",
			params: map[string]interface{}{
				"param1": 1,
			},
			expected: "SELECT * FROM table_1 WHERE status = ?",
		},
		{
			input: "SELECT * FROM table_1 WHERE status = :param1 and id = :param2",
			params: map[string]interface{}{
				"param1": 1,
				"param2": 2,
			},
			expected: "SELECT * FROM table_1 WHERE status = ? and id = ?",
		},
		{
			input: "SELECT * FROM table_1 WHERE status = :param1 and id = :param2 and amount = :param3",
			params: map[string]interface{}{
				"param1": 1,
				"param2": 2,
				"param3": 3,
			},
			expected: "SELECT * FROM table_1 WHERE status = ? and id = ? and amount = ?",
		},
	}

	db := &DB{driver: "mysql"}
	for _, v := range tests {
		result, _, err := db.BindNamed(v.input, v.params)
		if err != nil {
			t.Errorf("BindNamed() = %s, want not error", err)
		}
		if result != v.expected {
			t.Errorf("BindNamed() = %s, want %s", result, v.expected)
		}
	}
}

func Test_insertDriver(t *testing.T) {
	type usecase struct {
		input    string
		expected string
	}
	tests := []usecase{
		{
			input:    "postgres",
			expected: "postgres",
		},
		{
			input:    "mysql",
			expected: "mysql",
		},
		{
			input:    "nrpostgres",
			expected: "postgres",
		},
		{
			input:    "nrmysql",
			expected: "mysql",
		},
		{
			input:    "pgx",
			expected: "pgx",
		},
		{
			input:    "cloudsqlpostgres",
			expected: "cloudsqlpostgres",
		},
		{
			input:    "sqlite3",
			expected: "sqlite3",
		},
		{
			input:    "ora",
			expected: "ora",
		},
		{
			input:    "sqlserver",
			expected: "sqlserver",
		},
	}

	db := &DB{}
	for _, v := range tests {
		db.insertDriver(v.input)
		if db.driver != v.expected {
			t.Errorf("insertDriver() = %s, want %s", db.driver, v.expected)
		}
	}
}
