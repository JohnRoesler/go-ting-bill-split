Displays recent ting bills, split amongst lines

Install with go get
```bash
go get -u github.com/mastercactapus/go-ting-bill-split/tingbill
```

Run the command
```bash
tingbill
```

You will be prompted for your ting username and password. After fetching data it will print tables for the past 3 months with minutes, messages, and megabytes split up amongst the lines on the account.
```
+--------------+---------+----------+-----------+----------+--------+
|     NAME     | MINUTES | MESSAGES | MEGABYTES | BASE+TAX | TOTAL  |
+--------------+---------+----------+-----------+----------+--------+
| My Phone     | $1.52   | $3.61    | $5.87     | $7.76    | $18.76 |
| Other Phone  | $3.84   | $2.12    | $16.73    | $7.76    | $30.45 |
| Test Phone   | $3.64   | $2.27    | $6.39     | $7.76    | $20.06 |
|              |         |          |           |          | $69.27 |
+--------------+---------+----------+-----------+----------+--------+
```
