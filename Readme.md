TDS Commit Destroyer
===

##Usage

```
go get github.com/iwate/tds-commit-destroyer/cmd/tds-commit-destroyer
```

```
PS> tds-commit-destroyer -r {server-name}:1433
Proxying from :9433 to {server-name}:1433
```

Modify connection string like this
```
"Data Source=localhost,9433;Initial Catalog=sample;Integrated Security=True;Connect Timeout=2;Encrypt=False;TrustServerCertificate=False;ApplicationIntent=ReadWrite;MultiSubnetFailover=False"
```
`localhost,9433` is proxy endpoint. And you have to set `Encrypt` to `False`. 

##More 
TBD