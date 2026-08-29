# Changelog

## [1.1.0](https://github.com/peterintech/psocial/compare/v1.0.0...v1.1.0) (2026-08-29)


### Features

* update api version automatically ([7255d96](https://github.com/peterintech/psocial/commit/7255d96902c3a9ac56dbabd79885ca41ca53f1e8))

## 1.0.0 (2026-08-29)


### Features

* add basic auth and jwt authorization ([2816cd1](https://github.com/peterintech/psocial/commit/2816cd1c7e1edc735394198ac2de372461194f5d))
* add ci ([f03f69a](https://github.com/peterintech/psocial/commit/f03f69a2cfd73898ded8bb84ab12474540456609))
* add cors ([3865684](https://github.com/peterintech/psocial/commit/386568496e2bee49fe5ea46aacc6065df6f54985))
* add email serive with sendgrid to integrate invitation mail ([1e70a5d](https://github.com/peterintech/psocial/commit/1e70a5d87fd5523590fb6601cd7ce6313d7e86dc))
* add filter by search and tag to feed endpoint ([3ba6fcf](https://github.com/peterintech/psocial/commit/3ba6fcffe98b2a4372482ef176ce041089fb076a))
* add follow and unfollow apis ([8a0cf2e](https://github.com/peterintech/psocial/commit/8a0cf2e6edc24042266ad649d736b5a65a702816))
* add gracefull shutdown ([8757dec](https://github.com/peterintech/psocial/commit/8757dec43d969344734d1d5765d34c60c6e9375b))
* add indexes to optimize db search ([3415c8b](https://github.com/peterintech/psocial/commit/3415c8b516f2be89ee8c7b3d528f3aaa06884dc8))
* add metrics ([cdce8c2](https://github.com/peterintech/psocial/commit/cdce8c20bd6d2b060c02745fa0f57d56218134e4))
* add rate limiter using the fixed-window algo ([ae39692](https://github.com/peterintech/psocial/commit/ae39692c45d58fff75bb92977df46cacc1e0369e))
* add test to rate limiter ([28239ad](https://github.com/peterintech/psocial/commit/28239ad784d68ade1f67f81269121bb0d474ab7a))
* add token validation ([1aa2f53](https://github.com/peterintech/psocial/commit/1aa2f5396f6c0c93df0658bb369a6732c012eba7))
* add update and delete endpoints ([2ab4072](https://github.com/peterintech/psocial/commit/2ab407238ac6a9e0589fb884bb5ac638a6a9f06c))
* added caching with redis as a plugin ([667595c](https://github.com/peterintech/psocial/commit/667595c7d74636762fd1315a1e5e1340fd3d7cb3))
* **auth:** sql transactions for atomicity, add expires_at to invitation table and implement the logic ([a28ec90](https://github.com/peterintech/psocial/commit/a28ec90ed25d57d911565dd3735850271e0d4854))
* centralize error handling, add validator, add getpostby id api ([1c22632](https://github.com/peterintech/psocial/commit/1c2263246b22e36e3e6e042bcb212b6875d7f750))
* configure and run docker ([5652fe0](https://github.com/peterintech/psocial/commit/5652fe098db78e53f72e6f30c7e2ac4dddec69e8))
* database setup ([b7ef136](https://github.com/peterintech/psocial/commit/b7ef136c9fb4c176165924af946ded4860c14e24))
* DB migrated ([9cba192](https://github.com/peterintech/psocial/commit/9cba192153cc29f11317fa0a07f2d629bed0680d))
* get users ([3c08efa](https://github.com/peterintech/psocial/commit/3c08efa812f85ceca53168b0c9d6c9166bdeb1c5))
* implement create post/user ([9d8d09f](https://github.com/peterintech/psocial/commit/9d8d09f96db771e651a9309f96b22b2c3b5e23e3))
* implement RBAC to endpoints ([624b652](https://github.com/peterintech/psocial/commit/624b652b969d21bfcfb0e32e1e76efa300e61250))
* implement user feed algo, including pagination ([bb68aed](https://github.com/peterintech/psocial/commit/bb68aed300146eb2b20939f97650c52194691ff3))
* make a post endpoint done ([05ac5d8](https://github.com/peterintech/psocial/commit/05ac5d896c71dd8baad4b7d6a48a06a2f88f7d9f))
* migrate rbac authorizatiion to db ([767c405](https://github.com/peterintech/psocial/commit/767c405f128b2ce2c9a7660bc2b3d5cf9db2ed9d))
* optimistic concurrency control ([6c2d026](https://github.com/peterintech/psocial/commit/6c2d0265e4a4ebb04f8fd9a85df31d38702969c1))
* release pls action ([36b027c](https://github.com/peterintech/psocial/commit/36b027c18d75d290862b0462d9dfada56914c2e8))
* seed script, add comments api ([efe48d6](https://github.com/peterintech/psocial/commit/efe48d62b542cd3b20d6c3bc4d47f99fddde8c05))
* start auth setup ([66ffb5e](https://github.com/peterintech/psocial/commit/66ffb5ef5cd648b090dbc63abd8585e00a1edea8))
* store initialized ([8f0e007](https://github.com/peterintech/psocial/commit/8f0e007f9075f045ae4fb58d6e29ddc09360bd58))
* test concurrency, prevent racecondition ([6a9816f](https://github.com/peterintech/psocial/commit/6a9816fb534ec409e15bd5b7c1dd8506cd89517b))
* update seed, implement user invite endpoint ([1c34982](https://github.com/peterintech/psocial/commit/1c34982bdabe61f8ed8237c587a256a536755b2e))


### Bug Fixes

* add username to comments ([ae33432](https://github.com/peterintech/psocial/commit/ae3343229c2c6050e1ff81567205efc35ca63882))
* ensure correct password is provided to create token ([3bed229](https://github.com/peterintech/psocial/commit/3bed2290de3c57796d358d2f5d243cfb74ba86c3))
* point ci to master branch (base branch) ([74386ff](https://github.com/peterintech/psocial/commit/74386ff7613bb44e1b72a70c55d57ec0bcbb2bbe))
* ratelimiter test pass, bind host to ip for test, use contructor in test_util ([3a178ec](https://github.com/peterintech/psocial/commit/3a178ecf551eb9748ca829867f4dcc03858f332f))
