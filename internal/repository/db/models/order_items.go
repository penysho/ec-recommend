// Code generated by SQLBoiler 4.19.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// OrderItem is an object representing the database table.
type OrderItem struct {
	ID         string        `boil:"id" json:"id" toml:"id" yaml:"id"`
	OrderID    string        `boil:"order_id" json:"order_id" toml:"order_id" yaml:"order_id"`
	ProductID  string        `boil:"product_id" json:"product_id" toml:"product_id" yaml:"product_id"`
	Quantity   int           `boil:"quantity" json:"quantity" toml:"quantity" yaml:"quantity"`
	UnitPrice  types.Decimal `boil:"unit_price" json:"unit_price" toml:"unit_price" yaml:"unit_price"`
	TotalPrice types.Decimal `boil:"total_price" json:"total_price" toml:"total_price" yaml:"total_price"`
	CreatedAt  null.Time     `boil:"created_at" json:"created_at,omitempty" toml:"created_at" yaml:"created_at,omitempty"`

	R *orderItemR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L orderItemL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var OrderItemColumns = struct {
	ID         string
	OrderID    string
	ProductID  string
	Quantity   string
	UnitPrice  string
	TotalPrice string
	CreatedAt  string
}{
	ID:         "id",
	OrderID:    "order_id",
	ProductID:  "product_id",
	Quantity:   "quantity",
	UnitPrice:  "unit_price",
	TotalPrice: "total_price",
	CreatedAt:  "created_at",
}

var OrderItemTableColumns = struct {
	ID         string
	OrderID    string
	ProductID  string
	Quantity   string
	UnitPrice  string
	TotalPrice string
	CreatedAt  string
}{
	ID:         "order_items.id",
	OrderID:    "order_items.order_id",
	ProductID:  "order_items.product_id",
	Quantity:   "order_items.quantity",
	UnitPrice:  "order_items.unit_price",
	TotalPrice: "order_items.total_price",
	CreatedAt:  "order_items.created_at",
}

// Generated where

type whereHelpertypes_Decimal struct{ field string }

func (w whereHelpertypes_Decimal) EQ(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertypes_Decimal) NEQ(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertypes_Decimal) LT(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_Decimal) LTE(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_Decimal) GT(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_Decimal) GTE(x types.Decimal) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var OrderItemWhere = struct {
	ID         whereHelperstring
	OrderID    whereHelperstring
	ProductID  whereHelperstring
	Quantity   whereHelperint
	UnitPrice  whereHelpertypes_Decimal
	TotalPrice whereHelpertypes_Decimal
	CreatedAt  whereHelpernull_Time
}{
	ID:         whereHelperstring{field: "\"order_items\".\"id\""},
	OrderID:    whereHelperstring{field: "\"order_items\".\"order_id\""},
	ProductID:  whereHelperstring{field: "\"order_items\".\"product_id\""},
	Quantity:   whereHelperint{field: "\"order_items\".\"quantity\""},
	UnitPrice:  whereHelpertypes_Decimal{field: "\"order_items\".\"unit_price\""},
	TotalPrice: whereHelpertypes_Decimal{field: "\"order_items\".\"total_price\""},
	CreatedAt:  whereHelpernull_Time{field: "\"order_items\".\"created_at\""},
}

// OrderItemRels is where relationship names are stored.
var OrderItemRels = struct {
	Order   string
	Product string
}{
	Order:   "Order",
	Product: "Product",
}

// orderItemR is where relationships are stored.
type orderItemR struct {
	Order   *Order   `boil:"Order" json:"Order" toml:"Order" yaml:"Order"`
	Product *Product `boil:"Product" json:"Product" toml:"Product" yaml:"Product"`
}

// NewStruct creates a new relationship struct
func (*orderItemR) NewStruct() *orderItemR {
	return &orderItemR{}
}

func (o *OrderItem) GetOrder() *Order {
	if o == nil {
		return nil
	}

	return o.R.GetOrder()
}

func (r *orderItemR) GetOrder() *Order {
	if r == nil {
		return nil
	}

	return r.Order
}

func (o *OrderItem) GetProduct() *Product {
	if o == nil {
		return nil
	}

	return o.R.GetProduct()
}

func (r *orderItemR) GetProduct() *Product {
	if r == nil {
		return nil
	}

	return r.Product
}

// orderItemL is where Load methods for each relationship are stored.
type orderItemL struct{}

var (
	orderItemAllColumns            = []string{"id", "order_id", "product_id", "quantity", "unit_price", "total_price", "created_at"}
	orderItemColumnsWithoutDefault = []string{"order_id", "product_id", "quantity", "unit_price", "total_price"}
	orderItemColumnsWithDefault    = []string{"id", "created_at"}
	orderItemPrimaryKeyColumns     = []string{"id"}
	orderItemGeneratedColumns      = []string{}
)

type (
	// OrderItemSlice is an alias for a slice of pointers to OrderItem.
	// This should almost always be used instead of []OrderItem.
	OrderItemSlice []*OrderItem
	// OrderItemHook is the signature for custom OrderItem hook methods
	OrderItemHook func(context.Context, boil.ContextExecutor, *OrderItem) error

	orderItemQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	orderItemType                 = reflect.TypeOf(&OrderItem{})
	orderItemMapping              = queries.MakeStructMapping(orderItemType)
	orderItemPrimaryKeyMapping, _ = queries.BindMapping(orderItemType, orderItemMapping, orderItemPrimaryKeyColumns)
	orderItemInsertCacheMut       sync.RWMutex
	orderItemInsertCache          = make(map[string]insertCache)
	orderItemUpdateCacheMut       sync.RWMutex
	orderItemUpdateCache          = make(map[string]updateCache)
	orderItemUpsertCacheMut       sync.RWMutex
	orderItemUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var orderItemAfterSelectMu sync.Mutex
var orderItemAfterSelectHooks []OrderItemHook

var orderItemBeforeInsertMu sync.Mutex
var orderItemBeforeInsertHooks []OrderItemHook
var orderItemAfterInsertMu sync.Mutex
var orderItemAfterInsertHooks []OrderItemHook

var orderItemBeforeUpdateMu sync.Mutex
var orderItemBeforeUpdateHooks []OrderItemHook
var orderItemAfterUpdateMu sync.Mutex
var orderItemAfterUpdateHooks []OrderItemHook

var orderItemBeforeDeleteMu sync.Mutex
var orderItemBeforeDeleteHooks []OrderItemHook
var orderItemAfterDeleteMu sync.Mutex
var orderItemAfterDeleteHooks []OrderItemHook

var orderItemBeforeUpsertMu sync.Mutex
var orderItemBeforeUpsertHooks []OrderItemHook
var orderItemAfterUpsertMu sync.Mutex
var orderItemAfterUpsertHooks []OrderItemHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *OrderItem) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range orderItemAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *OrderItem) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range orderItemBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *OrderItem) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range orderItemAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *OrderItem) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range orderItemBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *OrderItem) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range orderItemAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *OrderItem) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range orderItemBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *OrderItem) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range orderItemAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *OrderItem) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range orderItemBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *OrderItem) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range orderItemAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddOrderItemHook registers your hook function for all future operations.
func AddOrderItemHook(hookPoint boil.HookPoint, orderItemHook OrderItemHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		orderItemAfterSelectMu.Lock()
		orderItemAfterSelectHooks = append(orderItemAfterSelectHooks, orderItemHook)
		orderItemAfterSelectMu.Unlock()
	case boil.BeforeInsertHook:
		orderItemBeforeInsertMu.Lock()
		orderItemBeforeInsertHooks = append(orderItemBeforeInsertHooks, orderItemHook)
		orderItemBeforeInsertMu.Unlock()
	case boil.AfterInsertHook:
		orderItemAfterInsertMu.Lock()
		orderItemAfterInsertHooks = append(orderItemAfterInsertHooks, orderItemHook)
		orderItemAfterInsertMu.Unlock()
	case boil.BeforeUpdateHook:
		orderItemBeforeUpdateMu.Lock()
		orderItemBeforeUpdateHooks = append(orderItemBeforeUpdateHooks, orderItemHook)
		orderItemBeforeUpdateMu.Unlock()
	case boil.AfterUpdateHook:
		orderItemAfterUpdateMu.Lock()
		orderItemAfterUpdateHooks = append(orderItemAfterUpdateHooks, orderItemHook)
		orderItemAfterUpdateMu.Unlock()
	case boil.BeforeDeleteHook:
		orderItemBeforeDeleteMu.Lock()
		orderItemBeforeDeleteHooks = append(orderItemBeforeDeleteHooks, orderItemHook)
		orderItemBeforeDeleteMu.Unlock()
	case boil.AfterDeleteHook:
		orderItemAfterDeleteMu.Lock()
		orderItemAfterDeleteHooks = append(orderItemAfterDeleteHooks, orderItemHook)
		orderItemAfterDeleteMu.Unlock()
	case boil.BeforeUpsertHook:
		orderItemBeforeUpsertMu.Lock()
		orderItemBeforeUpsertHooks = append(orderItemBeforeUpsertHooks, orderItemHook)
		orderItemBeforeUpsertMu.Unlock()
	case boil.AfterUpsertHook:
		orderItemAfterUpsertMu.Lock()
		orderItemAfterUpsertHooks = append(orderItemAfterUpsertHooks, orderItemHook)
		orderItemAfterUpsertMu.Unlock()
	}
}

// One returns a single orderItem record from the query.
func (q orderItemQuery) One(ctx context.Context, exec boil.ContextExecutor) (*OrderItem, error) {
	o := &OrderItem{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for order_items")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all OrderItem records from the query.
func (q orderItemQuery) All(ctx context.Context, exec boil.ContextExecutor) (OrderItemSlice, error) {
	var o []*OrderItem

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to OrderItem slice")
	}

	if len(orderItemAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all OrderItem records in the query.
func (q orderItemQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count order_items rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q orderItemQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if order_items exists")
	}

	return count > 0, nil
}

// Order pointed to by the foreign key.
func (o *OrderItem) Order(mods ...qm.QueryMod) orderQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.OrderID),
	}

	queryMods = append(queryMods, mods...)

	return Orders(queryMods...)
}

// Product pointed to by the foreign key.
func (o *OrderItem) Product(mods ...qm.QueryMod) productQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ProductID),
	}

	queryMods = append(queryMods, mods...)

	return Products(queryMods...)
}

// LoadOrder allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (orderItemL) LoadOrder(ctx context.Context, e boil.ContextExecutor, singular bool, maybeOrderItem interface{}, mods queries.Applicator) error {
	var slice []*OrderItem
	var object *OrderItem

	if singular {
		var ok bool
		object, ok = maybeOrderItem.(*OrderItem)
		if !ok {
			object = new(OrderItem)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeOrderItem)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeOrderItem))
			}
		}
	} else {
		s, ok := maybeOrderItem.(*[]*OrderItem)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeOrderItem)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeOrderItem))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &orderItemR{}
		}
		args[object.OrderID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &orderItemR{}
			}

			args[obj.OrderID] = struct{}{}

		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`orders`),
		qm.WhereIn(`orders.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Order")
	}

	var resultSlice []*Order
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Order")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for orders")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for orders")
	}

	if len(orderAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Order = foreign
		if foreign.R == nil {
			foreign.R = &orderR{}
		}
		foreign.R.OrderItems = append(foreign.R.OrderItems, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.OrderID == foreign.ID {
				local.R.Order = foreign
				if foreign.R == nil {
					foreign.R = &orderR{}
				}
				foreign.R.OrderItems = append(foreign.R.OrderItems, local)
				break
			}
		}
	}

	return nil
}

// LoadProduct allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (orderItemL) LoadProduct(ctx context.Context, e boil.ContextExecutor, singular bool, maybeOrderItem interface{}, mods queries.Applicator) error {
	var slice []*OrderItem
	var object *OrderItem

	if singular {
		var ok bool
		object, ok = maybeOrderItem.(*OrderItem)
		if !ok {
			object = new(OrderItem)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeOrderItem)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeOrderItem))
			}
		}
	} else {
		s, ok := maybeOrderItem.(*[]*OrderItem)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeOrderItem)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeOrderItem))
			}
		}
	}

	args := make(map[interface{}]struct{})
	if singular {
		if object.R == nil {
			object.R = &orderItemR{}
		}
		args[object.ProductID] = struct{}{}

	} else {
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &orderItemR{}
			}

			args[obj.ProductID] = struct{}{}

		}
	}

	if len(args) == 0 {
		return nil
	}

	argsSlice := make([]interface{}, len(args))
	i := 0
	for arg := range args {
		argsSlice[i] = arg
		i++
	}

	query := NewQuery(
		qm.From(`products`),
		qm.WhereIn(`products.id in ?`, argsSlice...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Product")
	}

	var resultSlice []*Product
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Product")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for products")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for products")
	}

	if len(productAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Product = foreign
		if foreign.R == nil {
			foreign.R = &productR{}
		}
		foreign.R.OrderItems = append(foreign.R.OrderItems, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ProductID == foreign.ID {
				local.R.Product = foreign
				if foreign.R == nil {
					foreign.R = &productR{}
				}
				foreign.R.OrderItems = append(foreign.R.OrderItems, local)
				break
			}
		}
	}

	return nil
}

// SetOrder of the orderItem to the related item.
// Sets o.R.Order to related.
// Adds o to related.R.OrderItems.
func (o *OrderItem) SetOrder(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Order) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"order_items\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"order_id"}),
		strmangle.WhereClause("\"", "\"", 2, orderItemPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.OrderID = related.ID
	if o.R == nil {
		o.R = &orderItemR{
			Order: related,
		}
	} else {
		o.R.Order = related
	}

	if related.R == nil {
		related.R = &orderR{
			OrderItems: OrderItemSlice{o},
		}
	} else {
		related.R.OrderItems = append(related.R.OrderItems, o)
	}

	return nil
}

// SetProduct of the orderItem to the related item.
// Sets o.R.Product to related.
// Adds o to related.R.OrderItems.
func (o *OrderItem) SetProduct(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Product) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"order_items\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"product_id"}),
		strmangle.WhereClause("\"", "\"", 2, orderItemPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ProductID = related.ID
	if o.R == nil {
		o.R = &orderItemR{
			Product: related,
		}
	} else {
		o.R.Product = related
	}

	if related.R == nil {
		related.R = &productR{
			OrderItems: OrderItemSlice{o},
		}
	} else {
		related.R.OrderItems = append(related.R.OrderItems, o)
	}

	return nil
}

// OrderItems retrieves all the records using an executor.
func OrderItems(mods ...qm.QueryMod) orderItemQuery {
	mods = append(mods, qm.From("\"order_items\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"order_items\".*"})
	}

	return orderItemQuery{q}
}

// FindOrderItem retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindOrderItem(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*OrderItem, error) {
	orderItemObj := &OrderItem{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"order_items\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, orderItemObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from order_items")
	}

	if err = orderItemObj.doAfterSelectHooks(ctx, exec); err != nil {
		return orderItemObj, err
	}

	return orderItemObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *OrderItem) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no order_items provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(orderItemColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	orderItemInsertCacheMut.RLock()
	cache, cached := orderItemInsertCache[key]
	orderItemInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			orderItemAllColumns,
			orderItemColumnsWithDefault,
			orderItemColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(orderItemType, orderItemMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(orderItemType, orderItemMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"order_items\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"order_items\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into order_items")
	}

	if !cached {
		orderItemInsertCacheMut.Lock()
		orderItemInsertCache[key] = cache
		orderItemInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the OrderItem.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *OrderItem) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	orderItemUpdateCacheMut.RLock()
	cache, cached := orderItemUpdateCache[key]
	orderItemUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			orderItemAllColumns,
			orderItemPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update order_items, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"order_items\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, orderItemPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(orderItemType, orderItemMapping, append(wl, orderItemPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update order_items row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for order_items")
	}

	if !cached {
		orderItemUpdateCacheMut.Lock()
		orderItemUpdateCache[key] = cache
		orderItemUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q orderItemQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for order_items")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for order_items")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o OrderItemSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), orderItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"order_items\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, orderItemPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in orderItem slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all orderItem")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *OrderItem) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("models: no order_items provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if queries.MustTime(o.CreatedAt).IsZero() {
			queries.SetScanner(&o.CreatedAt, currTime)
		}
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(orderItemColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	orderItemUpsertCacheMut.RLock()
	cache, cached := orderItemUpsertCache[key]
	orderItemUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			orderItemAllColumns,
			orderItemColumnsWithDefault,
			orderItemColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			orderItemAllColumns,
			orderItemPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert order_items, could not build update column list")
		}

		ret := strmangle.SetComplement(orderItemAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(orderItemPrimaryKeyColumns) == 0 {
				return errors.New("models: unable to upsert order_items, could not build conflict column list")
			}

			conflict = make([]string, len(orderItemPrimaryKeyColumns))
			copy(conflict, orderItemPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"order_items\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(orderItemType, orderItemMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(orderItemType, orderItemMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert order_items")
	}

	if !cached {
		orderItemUpsertCacheMut.Lock()
		orderItemUpsertCache[key] = cache
		orderItemUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single OrderItem record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *OrderItem) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no OrderItem provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), orderItemPrimaryKeyMapping)
	sql := "DELETE FROM \"order_items\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from order_items")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for order_items")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q orderItemQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no orderItemQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from order_items")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for order_items")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o OrderItemSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(orderItemBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), orderItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"order_items\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, orderItemPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from orderItem slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for order_items")
	}

	if len(orderItemAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *OrderItem) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindOrderItem(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *OrderItemSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := OrderItemSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), orderItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"order_items\".* FROM \"order_items\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, orderItemPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in OrderItemSlice")
	}

	*o = slice

	return nil
}

// OrderItemExists checks if the OrderItem row exists.
func OrderItemExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"order_items\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if order_items exists")
	}

	return exists, nil
}

// Exists checks if the OrderItem row exists.
func (o *OrderItem) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return OrderItemExists(ctx, exec, o.ID)
}
