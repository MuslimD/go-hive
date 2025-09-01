package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

func setupPostgres(t *testing.T) ParcelService {
	dsn := "muslimD:qwe12345@tcp(127.0.0.1:3406)/go_bd?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)

	require.NoError(t, err)
	require.NoError(t, db.Ping())

	store := NewParcelStore(db)
	service := NewParcelService(store)

	return service
}

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	service := setupPostgres(t)
	parcel := getTestParcel()

	id, err := service.store.Add(parcel)
	require.NoError(t, err)
	require.NotEqual(t, 0, id, "id не должен быть равен 0")

	p, err := service.store.Get(id)
	require.NoError(t, err)

	p.Number = 0
	p.CreatedAt = ""
	parcel.CreatedAt = ""

	require.Equal(t, p, parcel, "p должен быть равен parcel")

	require.NoError(t, service.store.Delete(id))
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {

	service := setupPostgres(t)
	parcel := getTestParcel()

	id, err := service.store.Add(parcel)
	require.NoError(t, err)
	require.NotEqual(t, 0, id, "id не должен быть равен 0")

	newAddress := "new test address"
	require.NoError(t, service.store.SetAddress(id, newAddress))

	p, err := service.store.Get(id)
	require.NoError(t, err)
	require.Equal(t, p.Address, newAddress, "адрес в новом parcel должен быть равен newAddress")
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	service := setupPostgres(t)
	parcel := getTestParcel()

	id, err := service.store.Add(parcel)
	require.NoError(t, err)
	require.NotEqual(t, 0, id, "id не должен быть равен 0")

	newStatus := "new test status"
	require.NoError(t, service.store.SetStatus(id, newStatus))

	p, err := service.store.Get(id)
	require.NoError(t, err)
	require.Equal(t, p.Status, newStatus, "status в новом parcel должен быть равен newStatus")
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	service := setupPostgres(t)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := service.store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEqual(t, 0, id, "id не должен быть равен 0")

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// get by client

	storedParcels, err := service.store.GetByClient(client)
	require.NoError(t, err)
	require.Equal(t, len(storedParcels), len(parcels), "количество полученных посылок должно быть равно количеству добавленных")

	// получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных

	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap

		p, ok := parcelMap[parcel.Number]
		require.True(t, ok, "ключ должен существовать в map")
		p.CreatedAt = ""
		parcel.CreatedAt = ""
		require.Equal(t, p, parcel)
		// убедитесь, что значения полей полученных посылок заполнены верно
	}
}
