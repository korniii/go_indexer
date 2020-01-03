// Code generated by SQLBoiler 3.6.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Customers", testCustomers)
	t.Run("Items", testItems)
	t.Run("Orders", testOrders)
}

func TestDelete(t *testing.T) {
	t.Run("Customers", testCustomersDelete)
	t.Run("Items", testItemsDelete)
	t.Run("Orders", testOrdersDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Customers", testCustomersQueryDeleteAll)
	t.Run("Items", testItemsQueryDeleteAll)
	t.Run("Orders", testOrdersQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Customers", testCustomersSliceDeleteAll)
	t.Run("Items", testItemsSliceDeleteAll)
	t.Run("Orders", testOrdersSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Customers", testCustomersExists)
	t.Run("Items", testItemsExists)
	t.Run("Orders", testOrdersExists)
}

func TestFind(t *testing.T) {
	t.Run("Customers", testCustomersFind)
	t.Run("Items", testItemsFind)
	t.Run("Orders", testOrdersFind)
}

func TestBind(t *testing.T) {
	t.Run("Customers", testCustomersBind)
	t.Run("Items", testItemsBind)
	t.Run("Orders", testOrdersBind)
}

func TestOne(t *testing.T) {
	t.Run("Customers", testCustomersOne)
	t.Run("Items", testItemsOne)
	t.Run("Orders", testOrdersOne)
}

func TestAll(t *testing.T) {
	t.Run("Customers", testCustomersAll)
	t.Run("Items", testItemsAll)
	t.Run("Orders", testOrdersAll)
}

func TestCount(t *testing.T) {
	t.Run("Customers", testCustomersCount)
	t.Run("Items", testItemsCount)
	t.Run("Orders", testOrdersCount)
}

func TestHooks(t *testing.T) {
	t.Run("Customers", testCustomersHooks)
	t.Run("Items", testItemsHooks)
	t.Run("Orders", testOrdersHooks)
}

func TestInsert(t *testing.T) {
	t.Run("Customers", testCustomersInsert)
	t.Run("Customers", testCustomersInsertWhitelist)
	t.Run("Items", testItemsInsert)
	t.Run("Items", testItemsInsertWhitelist)
	t.Run("Orders", testOrdersInsert)
	t.Run("Orders", testOrdersInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("ItemToOrderUsingOrder", testItemToOneOrderUsingOrder)
	t.Run("OrderToCustomerUsingCustomer", testOrderToOneCustomerUsingCustomer)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("CustomerToOrders", testCustomerToManyOrders)
	t.Run("OrderToItems", testOrderToManyItems)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("ItemToOrderUsingItems", testItemToOneSetOpOrderUsingOrder)
	t.Run("OrderToCustomerUsingOrders", testOrderToOneSetOpCustomerUsingCustomer)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("ItemToOrderUsingItems", testItemToOneRemoveOpOrderUsingOrder)
	t.Run("OrderToCustomerUsingOrders", testOrderToOneRemoveOpCustomerUsingCustomer)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("CustomerToOrders", testCustomerToManyAddOpOrders)
	t.Run("OrderToItems", testOrderToManyAddOpItems)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("CustomerToOrders", testCustomerToManySetOpOrders)
	t.Run("OrderToItems", testOrderToManySetOpItems)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("CustomerToOrders", testCustomerToManyRemoveOpOrders)
	t.Run("OrderToItems", testOrderToManyRemoveOpItems)
}

func TestReload(t *testing.T) {
	t.Run("Customers", testCustomersReload)
	t.Run("Items", testItemsReload)
	t.Run("Orders", testOrdersReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Customers", testCustomersReloadAll)
	t.Run("Items", testItemsReloadAll)
	t.Run("Orders", testOrdersReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Customers", testCustomersSelect)
	t.Run("Items", testItemsSelect)
	t.Run("Orders", testOrdersSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Customers", testCustomersUpdate)
	t.Run("Items", testItemsUpdate)
	t.Run("Orders", testOrdersUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Customers", testCustomersSliceUpdateAll)
	t.Run("Items", testItemsSliceUpdateAll)
	t.Run("Orders", testOrdersSliceUpdateAll)
}
