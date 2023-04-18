package jdb

import (
	"fmt"
	testing "testing"
)

func TestDatabaseCount(t *testing.T) {

	db, err := BuildDatabase("./db.json", encTodoData, decTodoData)
	if err != nil {
		t.Errorf("failed to read file ./db.json %v", err)
		return
	}

	ototal, err := db.CountItems()
	if err != nil {
		t.Errorf("failed to count original items %v", err)
		return
	}

	t.Logf("There was: %d\n", ototal)

	sd := TodoData{Name: fmt.Sprintf("Another thing to do: %d", (ototal + 1))}
	sd.ID = sd.NewID()

	saved, err := db.SaveItem(sd)
	if err != nil {
		t.Errorf("failed to saveitem %v with error %v", sd, err)
		return
	}

	t.Logf("Saved %v\n", saved)

	newtotal, err := db.CountItems()
	if err != nil {
		t.Errorf("failed to count new items %v", err)
		return
	}
	t.Logf("There is now: %d\n", newtotal)

	if ototal > newtotal {
		t.Errorf("new count %d was not larger than old count %d", newtotal, ototal)
		return
	}
}

func TestDatabaseGetItems(t *testing.T) {

	db, err := BuildDatabase("./db.json", encTodoData, decTodoData)
	if err != nil {
		t.Errorf("failed to read file ./db.json %v", err)
		return
	}

	sd := TodoData{Name: "at least one thing"}
	sd.ID = sd.NewID()

	_, err = db.SaveItem(sd)
	if err != nil {
		t.Errorf("failed to saveitem %v with error %v", sd, err)
		return
	}

	items, err := db.GetItems()
	if err != nil {
		t.Errorf("failed to get items %v", err)
		return
	}

	if len(items) == 0 {
		t.Error("expected at least 1 item in the database")
		return
	}

	for k, v := range items {
		t.Logf("%d: %v\n", k, v)
	}
}

func TestDatabaseGetItem(t *testing.T) {

	db, err := BuildDatabase("./db.json", encTodoData, decTodoData)
	if err != nil {
		t.Errorf("failed to read file ./db.json %v", err)
		return
	}

	idtotest := "b1f0e569-17ee-4d37-8b34-b35d9dd300cc"
	sd := TodoData{Name: "This has been updated"}
	sd.ID = idtotest

	_, err = db.SaveItem(sd)
	if err != nil {
		t.Errorf("failed to saveitem %v with error %v", sd, err)
		return
	}

	keyeditem, err := db.GetItem(idtotest)
	if err != nil {
		t.Errorf("failed to get items %v", err)
		return
	}
	readitem := keyeditem.(TodoData)

	if readitem.ID != idtotest {
		t.Errorf("new ID %s does not match read ID %s", readitem.ID, idtotest)
		return
	}
}

func TestDatabaseUpdate(t *testing.T) {

	db, err := BuildDatabase("./db.json", encTodoData, decTodoData)
	if err != nil {
		t.Errorf("failed to read file ./db.json %v", err)
		return
	}
	idtotest := "b1f0e569-17ee-4d37-8b34-b35d9dd300af"
	sd := TodoData{Name: "This has been updated"}
	sd.ID = idtotest

	saved, err := db.SaveItem(sd)
	if err != nil {
		t.Errorf("failed to saveitem %v with error %v", sd, err)
		return
	}

	t.Logf("Saved %v\n", saved)

	keyeditem, err := db.GetItem(idtotest)
	if err != nil {
		t.Errorf("failed to read item %s with err %v", idtotest, err)
		return
	}

	updateditem := keyeditem.(TodoData)

	t.Logf("updateditem is now: %v\n", updateditem)

	if updateditem.ID != sd.ID {
		t.Errorf("new ID %s does not match read ID %s", updateditem.ID, sd.ID)
		return
	}

	if updateditem.Name != sd.Name {
		t.Errorf("new Name %s does not match read Name %s", updateditem.Name, sd.Name)
		return
	}
}

func TestDatabaseDelete(t *testing.T) {

	db, err := BuildDatabase("./db.json", encTodoData, decTodoData)
	if err != nil {
		t.Errorf("failed to read file ./db.json %v", err)
		return
	}
	idtotest := "b1f0e569-17ee-4d37-8b34-b35d9dd300bb"
	sd := TodoData{Name: "This has been created or updated"}
	sd.ID = idtotest

	saved, err := db.SaveItem(sd)
	if err != nil {
		t.Errorf("failed to saveitem %v with error %v", sd, err)
		return
	}

	t.Logf("Saved %v\n", saved)

	keyeditem, err := db.GetItem(idtotest)
	if err != nil {
		t.Errorf("failed to read item %s with err %v", idtotest, err)
		return
	}

	updateditem := keyeditem.(TodoData)
	t.Logf("updateditem is now: %v\n", updateditem)

	err = db.DeleteItem(idtotest)
	if err != nil {
		t.Errorf("failed to delete item %s with err %v", idtotest, err)
		return
	}

	shouldnotexist, err := db.GetItem(idtotest)
	if err == nil {
		t.Errorf("read keyitem %s with data %v when we should not have", idtotest, shouldnotexist)
		return
	}

	t.Logf("got an error should be notexists: %v\n", err)

}
