package seeders

import (
	time "time"

	config "github.com/herumitra/ziidaapi/config"
	models "github.com/herumitra/ziidaapi/models"
)

func BranchSeed() {
	t := time.Date(2026, time.December, 31, 0, 0, 0, 0, time.UTC)

	branch := []models.Branch{
		{ID: "BRC11122024001", BranchName: "Branch 1", Address: "Jl. Raya Gudo, No. 101A, Gudo - Jombang", Phone: "085236990001", Email: "heru.oktafian@yahoo.com", SiaId: "xxxxxx", SiaName: "Mitra Farma", PsaId: "3517011710880001", PsaName: "Heru Oktafian", Sipa: "xxxxxx", SipaName: "Vita Fauzi. M", JournalMethod: "automatic", TaxPercentage: 11, BranchStatus: "active", LicenseDate: t},
		{ID: "BRC11122024002", BranchName: "Branch 2", Address: "Perum. Griya Nagari Singosari, Blok L, No. 1, Singosari - Malang", Phone: "0882009990001", Email: "oktafianheru@gmail.com", SiaId: "xxxxxx", SiaName: "Mitra Farma", PsaId: "3517011710880001", PsaName: "Heru Oktafian", Sipa: "xxxxxx", SipaName: "Vita Fauzi. M", JournalMethod: "automatic", TaxPercentage: 11, BranchStatus: "active", LicenseDate: t},
	}
	config.DB.Create(&branch)
}
