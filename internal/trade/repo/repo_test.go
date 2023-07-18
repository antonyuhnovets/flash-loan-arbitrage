package repo

import (
	"context"
	"os"
	"testing"

	. "github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	t "github.com/antonyuhnovets/flash-loan-arbitrage/internal/trade"
)

var basePath, _ = os.Getwd()
var testStorageLocal, _ = NewStorage(
	map[string]string{
		"pools":  string(basePath[:len(basePath)-23]) + "/storage_test/pools_test.json",
		"tokens": string(basePath[:len(basePath)-23]) + "/storage_test/tokens_test.json",
	},
)

var testStorages = []t.Repository{testStorageLocal}

func TestStorageSetup(t *testing.T) {
	for k, v := range testStorageLocal.fst.Files {
		_, err := os.Stat(v)
		if err != nil {
			err = testStorageLocal.fst.NewFile(k, v)
			if err != nil {
				t.Errorf(
					"recieved %s output on local setupStorageTest with index %v",
					err, k,
				)
			}
		} else {
			testStorageLocal.fst.UseFile(k, v)
			continue
		}
	}
}

type storeTest struct {
	where    string
	unit     interface{}
	expected error
}

var storeOneTests = []storeTest{
	{
		where: "tokens",
		unit: Token{
			Name:    "testOneTokenName",
			Address: "testOneTokenAddress",
			Wei:     10000000000,
		},
		expected: nil,
	},
	{
		where: "pools",
		unit: Pool{
			Address: "testOnePoolAddress",
			Pair: TokenPair{
				Token0: Token{Address: "testOnePairAddress0"},
				Token1: Token{Address: "testOnePairAddress1"},
			},
			Protocol: SwapProtocol{Name: "testOneProtocolName"},
		},
		expected: nil,
	},
}

func TestStoreOne(t *testing.T) {
	ctx := context.Background()

	for n, s := range testStorages {
		for index, el := range storeOneTests {
			switch el.where {
			case "tokens":
				err := s.AddToken(ctx, el.where, el.unit.(Token))
				if err == el.expected {
					continue
				} else {
					t.Errorf(
						"%s output on %v storage, %v TestStoreOne token",
						err, n, index,
					)
				}
			case "pools":
				err := s.AddPool(ctx, el.unit.(Pool), el.where)
				if err == el.expected {
					continue
				} else {
					t.Errorf(
						"%s output on %v storage, %v TestStoreOne pool",
						err, n, index,
					)
				}
			default:
				err := s.GetStorage().Store(ctx, el.where, el.unit.([]byte))
				if err == el.expected {
					continue
				} else {
					t.Errorf(
						"%s output on %v storage, %v TestStoreOne unknown",
						err, n, index,
					)
				}
			}
		}
	}
}

var storeManyTests = []storeTest{
	{
		where: "tokens",
		unit: []Token{
			{
				ID:      1,
				Name:    "testManyTokenName0",
				Address: "testManyTokenAddress0",
				Wei:     10000000000,
			}, {
				ID:      2,
				Name:    "testManyTokenName1",
				Address: "testManyTokenAddress1",
				Wei:     10000000000,
			},
		},
		expected: nil,
	},
	{
		where: "pools",
		unit: []Pool{
			{
				Address: "testManyPoolAddress0",
				Pair: TokenPair{
					Token0: Token{Address: "testManyPoolAddress0_0"},
					Token1: Token{Address: "testManyPoolAddress0_1"},
				},
				Protocol: SwapProtocol{Name: "testManyProtocolName0"},
			},
			{
				Address: "testManyPoolAddress1",
				Pair: TokenPair{
					Token0: Token{Address: "testManyPoolAddress1_0"},
					Token1: Token{Address: "testManyPoolAddress1_1"},
				},
				Protocol: SwapProtocol{Name: "testManyProtocolName0"},
			},
		},
		expected: nil,
	},
}

func TestStoreMany(t *testing.T) {
	ctx := context.Background()

	for n, s := range testStorages {
		for index, el := range storeManyTests {
			switch el.where {
			case "tokens":
				err := s.StoreTokens(ctx, el.where, el.unit.([]Token))
				if err == el.expected {
					continue
				} else {
					t.Errorf(
						"%s output on %v storage, %v TestStoreMany tokens",
						err, n, index,
					)
				}
			case "pools":
				err := s.StorePools(ctx, el.where, el.unit.([]Pool))
				if err == el.expected {
					continue
				} else {
					t.Errorf(
						"%s output on %v storage, %v TestStoreMany pools",
						err, n, index,
					)
				}
			default:
				err := s.GetStorage().Store(ctx, el.where, el.unit.([]byte))
				if err == el.expected {
					continue
				} else {
					t.Errorf(
						"%s output on %v storage, %v TestStoreMany unknown",
						err, n, index,
					)
				}
			}
		}
	}
}

type readTest struct {
	where    string
	key      interface{}
	expected interface{}
}

var getByTests = []readTest{
	{
		where:    "tokens",
		key:      "testManyTokenAddress1",
		expected: nil,
	},
	{
		where: "pools",
		key: TokenPair{
			Token0: Token{Address: "testManyPoolAddress0_0"},
			Token1: Token{Address: "testManyPoolAddress0_1"},
		},
		expected: nil,
	},
}

func TestGetBy(t *testing.T) {
	ctx := context.Background()

	for n, storage := range testStorages {
		for index, caseRT := range getByTests {
			switch caseRT.where {
			case "tokens":
				out, err := storage.GetTokenByAddress(
					ctx,
					caseRT.where, caseRT.key.(string),
				)
				if err == caseRT.expected {
					continue
				} else {
					t.Errorf(
						"%v %s output on %v storage, %v TestReadBy tokens",
						out, err, n, index,
					)
				}
			case "pools":
				pair, ok := caseRT.key.(TokenPair)
				if !ok {
					t.Errorf(
						"invalid key %v on %d storage, %d TestReadBy pools",
						caseRT.key, n, index,
					)
				}
				out, err := storage.GetByTokens(ctx, caseRT.where, pair)
				if err == caseRT.expected {
					continue
				} else {
					t.Errorf(
						"%v %s  on %v storage, %v TestReadBy pools",
						out, err, n, index,
					)
				}
			}
		}
	}
}

var readAllTests = []readTest{
	{
		where:    "tokens",
		expected: nil,
	},
	{
		where:    "pools",
		expected: nil,
	},
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()

	for n, s := range testStorages {
		for index, el := range readAllTests {
			switch el.where {
			case "tokens":
				lst, err := s.ListTokens(ctx, el.where)
				if err == el.expected {
					continue
				} else {
					t.Errorf(
						"%v %s output on %v storage, %v TestReadAll tokens",
						lst, err, n, index,
					)
				}
			case "pools":
				lst, err := s.ListPools(ctx, el.where)
				if err == el.expected {
					continue
				} else {
					t.Errorf(
						"%v %s output on %v storage, %v TestReadAll pools",
						lst, err, n, index,
					)
				}
			default:
				lst := []Pool{}
				err := s.GetStorage().Read(ctx, el.where, lst)
				if &lst == el.expected || err == el.expected {
					continue
				} else {
					t.Errorf(
						"%v %s output on %v storage, %v TestReadAll unknown",
						lst, err, n, index,
					)
				}
			}
		}
	}
}

func TestClean(t *testing.T) {
	ctx := context.Background()

	for index, storage := range testStorages {
		err := storage.GetStorage().ClearAll(ctx)
		if err != nil {
			t.Errorf(
				"error %s catched while cleaning storage %v",
				err, index,
			)
		}
	}
}

type RepoTestcase struct {
	tests map[string]func(*testing.T)
	order []string
}

var repoTestCases = []RepoTestcase{
	{tests: map[string]func(t *testing.T){
		"clean":        TestClean,
		"setupStorage": TestStorageSetup,
		"storeOne":     TestStoreOne,
		"storeMany":    TestStoreMany,
		"getBy":        TestGetBy,
		"getAll":       TestGetAll,
	},
		order: []string{
			"setupStorage",
			"storeOne",
			"storeMany",
			"getAll",
			"getBy",
			"clean",
		},
	},
}

func TestRepo(t *testing.T) {
	for index, tcase := range repoTestCases {
		for _, key := range tcase.order {
			t.Logf("testing %v", key)
			tcase.tests[key](t)
		}
		t.Logf("Testcase %v completed", index)
	}
}

// func _sameLists(list1, list2 []interface{}) bool {
// 	for index, el := range list1 {
// 		if list2[index] == el {
// 			continue
// 		} else {
// 			return false
// 		}
// 	}
// 	return true

// }
