payme/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── pkg/
│   ├── auth/                 # Authentication logic
│   ├── wallet/               # Wallet management
│   ├── transaction/          # Transactions and logs
│   ├── transfer/             # Internal transfers and funding
│   ├── savings/              # Group and Personal savings
│   ├── utilities/           # Airtime and Electricity
│   ├── utils/                # General helpers
│   └── database/             # DB connection logic
├── migrations/               # SQL migration files
│   └── 000001_initial_schema.up.sql
├── go.mod
└── go.sum
