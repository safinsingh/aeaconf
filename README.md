# aeaconf

DSL for aeacus based on INI format more intuitive boolean expressions for complex checks. Implements recursive descent parsing. Has support for parameterized custom conditions (see `testing/scoring.ini`).

## terminology

check: {message, points, complex condition}
condition: any complex boolean expression
func: PathExists, etc

## examples

`go run .` will transform `testing/scoring.ini` into

```
(main.Config) {
 Round: (main.Round) {
  Title: (string) (len=9) "Linux ICC",
  Os: (string) (len=15) "Ubuntu 20.04.03",
  User: (string) (len=7) "cpadmin",
  Local: (string) (len=5) "false"
 },
 Remote: (main.Remote) {
  Enable: (bool) true,
  Name: (bool) false,
  Server: (bool) false,
  Password: (bool) false
 },
 CustomConditions: (map[string]main.Condition) <nil>,
 Checks: ([]main.Check) (len=2 cap=2) {
  (main.Check) {
   Message: (string) (len=23) "Pam Password File works",
   Points: (int) 1,
   Cond: (main.OrExpr) {
    Lhs: (main.AndExpr) {
     Lhs: (main.PathExists) {
      Path: (string) (len=11) "/etc/passwd"
     },
     Rhs: (main.FileContains) {
      File: (string) (len=14) "/root/password",
      Value: (string) (len=5) "hello"
     }
    },
    Rhs: (main.PathExists) {
     Path: (string) (len=15) "/etc/pam.d/sshd"
    }
   }
  },
  (main.Check) {
   Message: (string) (len=20) "Forensics Question 1",
   Points: (int) 2,
   Cond: (main.FileContains) {
    File: (string) (len=13) "/home/desktop",
    Value: (string) (len=4) "abcd"
   }
  }
 }
}
```
