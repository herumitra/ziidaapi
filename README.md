## Tentang Ziida API
`(Golang Fiber Restful API with JWT)`
Dalam repositori ini kita menerapkan `Golang` sebagai platform dasar bahasa pemrograman yang digunakan dalam pembuatan `API`.
Di dalam repositori ini juga kami terapkan framework `Fiber` serta dependensi `GORM` dan `JWT` untuk mempermudah dalam pengerjaan di ranah sekuritas maupun pengelolaan databasenya, sehingga detail komponen yang kami gunakan bisa dijabarkan seperti berikut ini :
| NO. | KOMPONEN       |
|-----:|---------------|
|     1| Fiber         |
|     2| GORM          |
|     3| PostgreSQL    |
|     4| Redis         |
|     5| JWT           |


>Tambahkan file `.env` dalam direktori paling luar dari project di repositori ini, dan masukkan teks berikut di dalamnya.

```bash
JWT_SECRET=Sikrit1234
DB_USER=ziida
DB_PASSWORD=Pass1234
DB_NAME=ziida
DB_HOST=localhost
DB_PORT=5432
REDIS_HOST=localhost
REDIS_PORT=6379
SERVER_PORT=4001
```

>Jalankan perintah di bawah ini dalam query database `PostgreSQL`:S
```bash
CREATE TYPE journal_method AS ENUM ('manual', 'automatic');
CREATE TYPE user_role AS ENUM ('operator', 'administrator', 'cashier', 'finance');
CREATE TYPE status_user AS ENUM ('active', 'inactive');
```

>Jalankan perintah ini untuk generate tabel-tabel di database `PostgreSQL` serta seed data pengguna:
```bash
go run main.go seed
```