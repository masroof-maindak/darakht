---
tags:
  - passwd
  - C
  - getpwent
  - endpwent
  - getpwuid
  - getpwnam
site:
  - https://www.ibm.com/docs/en/aix/7.2?topic=g-getpwent-getpwuid-getpwnam-putpwent-setpwent-endpwent-subroutine
  - https://www.man7.org/linux/man-pages/man3/endpwent.3p.html
---

#### Fetching

Obtain the first matching entry by UserID, Name, or iterate through all of them sequentially by calling `getpwnam`, `getpwuid`, or `getpwent` respectively.

```C
struct passwd *getpwnam(const char *name);
struct passwd *getpwuid(uid_t uid);
struct passwd *getpwent(void);
```

#### Resetting & Closing

You can rewind to the start of the **[[passwd]]** file's entries with `setpwent()`, and when you're done, close the database with `endpwent()`.

```C
void setpwent(void);
void endpwent(void);
```

#### Struct

```C
#include <pwd.h>

struct passwd {
    char   *pw_name;       /* username */
    char   *pw_passwd;     /* user password */
    uid_t   pw_uid;        /* user ID */
    gid_t   pw_gid;        /* group ID */
    char   *pw_gecos;      /* user information */
    char   *pw_dir;        /* home directory */
    char   *pw_shell;      /* shell program */
};
```
