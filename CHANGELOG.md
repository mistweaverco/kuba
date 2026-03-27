## [1.8.3](https://github.com/mistweaverco/kuba/compare/v1.8.2...v1.8.3) (2026-03-27)


### Bug Fixes

* **ci:** update changelog generation script ([4477e66](https://github.com/mistweaverco/kuba/commit/4477e66c9b52dd20f4a1d23b1f4c1077c063fdd3))

## [1.8.2](https://github.com/mistweaverco/kuba/compare/v1.8.1...v1.8.2) (2026-03-26)


### Bug Fixes

* **ci:** try new apporach for changelog script ([a4cd79d](https://github.com/mistweaverco/kuba/commit/a4cd79dfe5107dc166573fb415582ce05c7d53de))

## [1.8.1](https://github.com/mistweaverco/kuba/compare/v1.8.0...v1.8.1) (2026-03-26)


### Bug Fixes

* **ci:** changelog generation script ([6e5ecb0](https://github.com/mistweaverco/kuba/commit/6e5ecb0a4b813116e27c7f1a14c80b4cb8b58037))

# [1.8.0](https://github.com/mistweaverco/kuba/compare/v1.7.0...v1.8.0) (2026-03-26)


### Bug Fixes

* **ci:** add depends on for build steps ([130d1f7](https://github.com/mistweaverco/kuba/commit/130d1f73e44c5d831e9a2aeaa86f2da3fdedda75))
* **ci:** remove changelog gen step from build script ([39f3d82](https://github.com/mistweaverco/kuba/commit/39f3d82028d7b357c5f9ada2f84702d48ea19ef8))
* **ci:** update changelog generation script ([046a658](https://github.com/mistweaverco/kuba/commit/046a6581adb18b7cc511e844ff1bc168fb7fbe48))
* **ci:** web ([9fd6808](https://github.com/mistweaverco/kuba/commit/9fd6808008f89ddd9bcbd5d257abe1a814223b6d))
* **gitignore:** Update .gitignore to exclude build artifacts ([bfeeb8c](https://github.com/mistweaverco/kuba/commit/bfeeb8c364f71e53ddaeaebce781b1515b73e88f))
* **tui:** overflows ([#68](https://github.com/mistweaverco/kuba/issues/68)) ([32292f5](https://github.com/mistweaverco/kuba/commit/32292f505ce344e864e0adfbd8f34918d555531e))


### Features

* **cli:** add kuba changelog command ([#71](https://github.com/mistweaverco/kuba/issues/71)) ([27b823a](https://github.com/mistweaverco/kuba/commit/27b823aa4f03c0eef6b517fd70e439f0c69925be))
* **cli:** add new `create template` command ([#70](https://github.com/mistweaverco/kuba/issues/70)) ([0089566](https://github.com/mistweaverco/kuba/commit/0089566842a3c11eaed33cdabb1d4eae88e1a6a0))
* **tui:** change keybinds ([#67](https://github.com/mistweaverco/kuba/issues/67)) ([1d13bda](https://github.com/mistweaverco/kuba/commit/1d13bda277e889525af4567d73bc0a0ef82927b6))

# [1.7.0](https://github.com/mistweaverco/kuba/compare/v1.6.2...v1.7.0) (2026-03-25)


### Bug Fixes

* **ci:** update go to 1.25 ([#65](https://github.com/mistweaverco/kuba/issues/65)) ([397b9c6](https://github.com/mistweaverco/kuba/commit/397b9c6b93815967c8813621c5f14a71c3aed865))


### Features

* **tui:** add tui for viewing, editing and adding secrets ([#64](https://github.com/mistweaverco/kuba/issues/64)) ([31d37bd](https://github.com/mistweaverco/kuba/commit/31d37bd2e6c09f56f1956638a8acabe0844644dd))
* **ux:** add changelog generation script ([#63](https://github.com/mistweaverco/kuba/issues/63)) ([17cb5d9](https://github.com/mistweaverco/kuba/commit/17cb5d9a31f5aeaf7484f9068435d154b838a06f))

## [1.6.2](https://github.com/mistweaverco/kuba/compare/v1.6.1...v1.6.2) (2026-03-24)


### Bug Fixes

* **windows:** crash due to missing bitwarden dll ([#59](https://github.com/mistweaverco/kuba/issues/59)) ([7a61b69](https://github.com/mistweaverco/kuba/commit/7a61b69a210ee8781cc929002d862097e0a25460))
* **windows:** crash due to wrong filepath join ([#60](https://github.com/mistweaverco/kuba/issues/60)) ([e3170a8](https://github.com/mistweaverco/kuba/commit/e3170a84a760b7527180bf1fb9cd67628ba17a42))
* **windows:** kuba update ([#58](https://github.com/mistweaverco/kuba/issues/58)) ([a005678](https://github.com/mistweaverco/kuba/commit/a00567859e3cbd96f4fc86cb23367d3d24680e46))

## [1.6.1](https://github.com/mistweaverco/kuba/compare/v1.6.0...v1.6.1) (2026-03-24)


### Bug Fixes

* **ci:** pkgname was missing -bin suffix ([#57](https://github.com/mistweaverco/kuba/issues/57)) ([cb8478c](https://github.com/mistweaverco/kuba/commit/cb8478c04ed3e9a744397163ba7e0a49c2ed8b3e))
* **windows:** crash on windows due to bitwarden sdk bug ([#55](https://github.com/mistweaverco/kuba/issues/55)) ([87ec150](https://github.com/mistweaverco/kuba/commit/87ec150a58f60a3cc28188fabf9e83b195277c78))
* **windows:** Don't use hardcoded temp directory ([#54](https://github.com/mistweaverco/kuba/issues/54)) ([23d42d7](https://github.com/mistweaverco/kuba/commit/23d42d7e5cf014f91765bc90824769fa641d24da))
* **windows:** install.ps1 caused syntax errors ([#56](https://github.com/mistweaverco/kuba/issues/56)) ([a95a9c1](https://github.com/mistweaverco/kuba/commit/a95a9c10c81a87939ed40e3b3f695291fd146c71))


### Features

* **CI:** Add PKGBUILD for Arch Linux + release workflow ([#53](https://github.com/mistweaverco/kuba/issues/53)) ([bbb8bc2](https://github.com/mistweaverco/kuba/commit/bbb8bc222961dce6ebdfa59046eeaace62d0f71e))

# [1.6.0](https://github.com/mistweaverco/kuba/compare/v1.5.0...v1.6.0) (2026-03-13)


### Features

* **convert:** extend from kscv to support remote imports ([#52](https://github.com/mistweaverco/kuba/issues/52)) ([e74a3c6](https://github.com/mistweaverco/kuba/commit/e74a3c66bbed871f6ad0d1372168360a7b703698))

# [1.5.0](https://github.com/mistweaverco/kuba/compare/v1.4.0...v1.5.0) (2026-03-02)


### Bug Fixes

* **ci:** add c dep for linux ([886fd37](https://github.com/mistweaverco/kuba/commit/886fd371fc4d1049eaf612314fb71daafd634863))
* **ci:** build for linux and windows ([6c1a4c8](https://github.com/mistweaverco/kuba/commit/6c1a4c8b7e94c505b8b96da4fea5787950b0547e))
* **ci:** build for windows needs a windows runner ([d3c424d](https://github.com/mistweaverco/kuba/commit/d3c424df60b2a746266dcc4cdbf115c51b45a45a))
* **ci:** conflict with windows reserved PLATFORM env var ([656afe7](https://github.com/mistweaverco/kuba/commit/656afe7ca26543fead13cbf2396c3b65dca1c123))
* **ci:** enable CGO for bitwarden ([09a3b13](https://github.com/mistweaverco/kuba/commit/09a3b13762ea43122a774feeaa3a28042f4857b3))
* **ci:** make cc available in path ([da5b9d0](https://github.com/mistweaverco/kuba/commit/da5b9d07884e7ce622d89339c516c5cf768a9c20))
* **ci:** release.sh still referenced PLATFORM instead of TARGET_* ([277a55f](https://github.com/mistweaverco/kuba/commit/277a55f862a651a2624fd75c904814c3246733b5))
* **ci:** remove 32 bit version ([45bc035](https://github.com/mistweaverco/kuba/commit/45bc035422f6d8ab19b67897d81e0e57f9e0e54c))
* **ci:** use correct c toolchain ([3bf7ace](https://github.com/mistweaverco/kuba/commit/3bf7ace583003269a3ec4b1c0b6e0532ccf68cf9))
* **ci:** use mingw gcc for cgo on windows ([39c6b26](https://github.com/mistweaverco/kuba/commit/39c6b2662b298084052c0ff0b386835af70a7c63))


### Features

* **docs:** add bitwarden to TOC iin README ([f5fdd76](https://github.com/mistweaverco/kuba/commit/f5fdd76a8d7366513fe9ca70dfd1ba43784284d1))
* **providers:** add bitwarden support ([#50](https://github.com/mistweaverco/kuba/issues/50)) ([7b7917a](https://github.com/mistweaverco/kuba/commit/7b7917a9fd75ce89142082f7e90563156f103b9a))

# [1.4.0](https://github.com/mistweaverco/kuba/compare/v1.3.0...v1.4.0) (2026-03-02)


### Features

* **convert:** add "import" from knative service (ksvc) ([#48](https://github.com/mistweaverco/kuba/issues/48)) ([abe6f1b](https://github.com/mistweaverco/kuba/commit/abe6f1b822025b860fd3a9bcfd316fbaef2504fc))
* **docs:** update docs for "import from knative service" ([#49](https://github.com/mistweaverco/kuba/issues/49)) ([41ef6f1](https://github.com/mistweaverco/kuba/commit/41ef6f10004e151ff1e90ed2206fde7244543aab))

# [1.3.0](https://github.com/mistweaverco/kuba/compare/v1.2.0...v1.3.0) (2025-12-15)


### Bug Fixes

* **web:** spelling ([71c153b](https://github.com/mistweaverco/kuba/commit/71c153b603e4ed37ae765a21dffbe6b5f1975a02))


### Features

* **cli:** `kuba show --env` lists all available envs now ([#46](https://github.com/mistweaverco/kuba/issues/46)) ([bce2942](https://github.com/mistweaverco/kuba/commit/bce2942c0ab1f32e654624625b052304c1c9c73f))
* **cli:** `kuba show` now support `-o`/`--output` flag ([#47](https://github.com/mistweaverco/kuba/issues/47)) ([a35a6ab](https://github.com/mistweaverco/kuba/commit/a35a6abd8ed4212516a4819d6d4062bcf830f0a3))

# [1.2.0](https://github.com/mistweaverco/kuba/compare/v1.1.0...v1.2.0) (2025-12-06)


### Bug Fixes

* **cli:** remove commands registered twice ([#40](https://github.com/mistweaverco/kuba/issues/40)) ([91631b5](https://github.com/mistweaverco/kuba/commit/91631b58220b5aaea46c639158420e11ea3a397d))
* **docs:** update README ([#39](https://github.com/mistweaverco/kuba/issues/39)) ([10fc57b](https://github.com/mistweaverco/kuba/commit/10fc57bfa23dd53758d97a6f8d5b7e6db3967448))


### Features

* **cli:** add `--command` flag to `run` command ([#41](https://github.com/mistweaverco/kuba/issues/41)) ([849d1a8](https://github.com/mistweaverco/kuba/commit/849d1a87283830d88fb531abb31845a2b97464c3))

# [1.1.0](https://github.com/mistweaverco/kuba/compare/v1.0.0...v1.1.0) (2025-12-05)


### Features

* **convert:** add convert command ([#35](https://github.com/mistweaverco/kuba/issues/35)) ([518bcac](https://github.com/mistweaverco/kuba/commit/518bcac8923cadaf705b59f3cd4203beb4ee2ff2))
* **dx:** add show command ([#37](https://github.com/mistweaverco/kuba/issues/37)) ([355b486](https://github.com/mistweaverco/kuba/commit/355b486785c166036f31d1c18ecbdd43d1d882f1))
* **test:** improve test command and give better feedback ([#36](https://github.com/mistweaverco/kuba/issues/36)) ([5aadad9](https://github.com/mistweaverco/kuba/commit/5aadad9a211da2eb067ed06f5ca838b627d962d9))

# [1.0.0](https://github.com/mistweaverco/kuba/compare/v0.10.0...v1.0.0) (2025-10-04)


### Features

* **cache:** add management commands ([#29](https://github.com/mistweaverco/kuba/issues/29)) ([11b361a](https://github.com/mistweaverco/kuba/commit/11b361a8639a9a7f8be777fab3df1c3c5f16d9bc))

# [0.10.0](https://github.com/mistweaverco/kuba/compare/v0.9.4...v0.10.0) (2025-10-04)


### Features

* **cache:** add cache ([#28](https://github.com/mistweaverco/kuba/issues/28)) ([bd8348b](https://github.com/mistweaverco/kuba/commit/bd8348bc6fb75ef1fcfff088b6a112feb2d37a70))
* **kuba:** add update subcommand ([#27](https://github.com/mistweaverco/kuba/issues/27)) ([f199aec](https://github.com/mistweaverco/kuba/commit/f199aec0472120fe719d3e2cfabe5d98941f80cd))

## [0.9.4](https://github.com/mistweaverco/kuba/compare/v0.9.3...v0.9.4) (2025-10-04)


### Bug Fixes

* **manager:** correctly interpolate variables ([#26](https://github.com/mistweaverco/kuba/issues/26)) ([312f705](https://github.com/mistweaverco/kuba/commit/312f705591a7c9c96911ab984b798275440bc3f8))


### Features

* **docs:** add development-status badge ([0669b33](https://github.com/mistweaverco/kuba/commit/0669b33e81748e6682a98300b6865d5386923aff))

## [0.9.3](https://github.com/mistweaverco/kuba/compare/v0.9.2...v0.9.3) (2025-09-19)


### Bug Fixes

* **config:** multiple var interpolation ([#22](https://github.com/mistweaverco/kuba/issues/22)) ([1fe6f23](https://github.com/mistweaverco/kuba/commit/1fe6f23041ec7a24edd1d5c1b921a99028ae7260))

## [0.9.2](https://github.com/mistweaverco/kuba/compare/v0.9.1...v0.9.2) (2025-09-18)


### Bug Fixes

* **kuba-init:** replace mappings with env key ([#21](https://github.com/mistweaverco/kuba/issues/21)) ([46c4171](https://github.com/mistweaverco/kuba/commit/46c4171d3f8620eb75e5f0730deb0c7bfd1b15f7))

## [0.9.1](https://github.com/mistweaverco/kuba/compare/v0.9.0...v0.9.1) (2025-09-18)


### Bug Fixes

* **init:** kuba init produced old kuba.yaml format ([#20](https://github.com/mistweaverco/kuba/issues/20)) ([0710f3f](https://github.com/mistweaverco/kuba/commit/0710f3fd0be6c8f782d63df6d4b2dd4c9a138783))
* **schema:** kuba.schema.json was enforcing wrong rules ([#19](https://github.com/mistweaverco/kuba/issues/19)) ([ade9c37](https://github.com/mistweaverco/kuba/commit/ade9c37dd8e7c425ba7d1319ed74de13b4fb5e0a))

# [0.9.0](https://github.com/mistweaverco/kuba/compare/v0.8.0...v0.9.0) (2025-09-18)


### Bug Fixes

* **docs:** docker --env-file needs = ([d1a6ddd](https://github.com/mistweaverco/kuba/commit/d1a6ddde3891c76ecb33b6196adbe397511d0591))
* **docs:** docker usage ([#14](https://github.com/mistweaverco/kuba/issues/14)) ([218ae12](https://github.com/mistweaverco/kuba/commit/218ae122f4c0c8f5a2f50cad67c1ad6d428d56a3))


### Features

* **config:** add inherits from env block ([#18](https://github.com/mistweaverco/kuba/issues/18)) ([96d1b7e](https://github.com/mistweaverco/kuba/commit/96d1b7e8bd99ca7dfcc86d011cd774296a8141f5))

# [0.8.0](https://github.com/mistweaverco/kuba/compare/v0.7.1...v0.8.0) (2025-09-11)


### Bug Fixes

* **web:** install.sh ([f85eb61](https://github.com/mistweaverco/kuba/commit/f85eb61b5f66947a8f34bbc8d2ace6e8f4c909eb))
* **web:** install.sh move echo ([594b1f7](https://github.com/mistweaverco/kuba/commit/594b1f7548269fa376ee5d68fc929d3ff4b04ece))


### Features

* **flag:** add contain flag ([#13](https://github.com/mistweaverco/kuba/issues/13)) ([99409bb](https://github.com/mistweaverco/kuba/commit/99409bb6535e4db99c28475a4cef09ca0792a9ee))

## [0.7.1](https://github.com/mistweaverco/kuba/compare/v0.7.0...v0.7.1) (2025-09-10)


### Bug Fixes

* **ci:** apple sign and notarize ([694d0cb](https://github.com/mistweaverco/kuba/commit/694d0cb052f921055ae8096951943a7e6c5f6273))
* **ci:** mac builds need to run macos ([4716fc8](https://github.com/mistweaverco/kuba/commit/4716fc8574437d7ed07d1e3d384c107e6d00609b))
* **ci:** sign and notarize for macos ([#12](https://github.com/mistweaverco/kuba/issues/12)) ([7c8a7dd](https://github.com/mistweaverco/kuba/commit/7c8a7dd79009d7d64175c9fe287a72d031c8f981))

# [0.7.0](https://github.com/mistweaverco/kuba/compare/v0.6.0...v0.7.0) (2025-09-10)


### Bug Fixes

* **ClickableHeadline:** make it more pleasant to view at ([f7c072a](https://github.com/mistweaverco/kuba/commit/f7c072a367f7753e03d1429681830e149d6d3215))
* **examples:** also use clickable headline component ([7e5fdb3](https://github.com/mistweaverco/kuba/commit/7e5fdb3f722c6c650deb14707b584d4f0503e33a))
* **web:** add json highlighting ([ea75bbe](https://github.com/mistweaverco/kuba/commit/ea75bbe2048061d3c457d33481f511a97a7b5360))
* **web:** fix colors ([8238bad](https://github.com/mistweaverco/kuba/commit/8238bada20734de1794332736754c2b127684b87))
* **web:** linter ([98e640b](https://github.com/mistweaverco/kuba/commit/98e640bfb2f78445cfbc6c6201c1a44034658f98))


### Features

* **cmd:** add test subcommand ([#11](https://github.com/mistweaverco/kuba/issues/11)) ([59cd9e7](https://github.com/mistweaverco/kuba/commit/59cd9e7d50d4b59a1654ced8fbab61e808671c0b))
* **web:** add clickable headlines ([25b838c](https://github.com/mistweaverco/kuba/commit/25b838ced7122c16aecd1c47f798fd7e0554ea1b))
* **web:** add open-graph ([f194805](https://github.com/mistweaverco/kuba/commit/f194805efaf3d9d50b181564960d2bfaec203b05))
* **web:** add toast to ClickableHeadline ([07e15dc](https://github.com/mistweaverco/kuba/commit/07e15dc4658d74071d9bb10493559be10ccc13bf))
* **web:** use dracula daisyui theme ([4d40c85](https://github.com/mistweaverco/kuba/commit/4d40c859a1c82c4c4005c765dc23841eecf7569b))

# [0.6.0](https://github.com/mistweaverco/kuba/compare/v0.5.0...v0.6.0) (2025-08-25)


### Bug Fixes

* **web:** build ([843050c](https://github.com/mistweaverco/kuba/commit/843050cd158049f6b58ccce696291f211ec24b50))


### Features

* **cmd:** add debug mode flag ([#10](https://github.com/mistweaverco/kuba/issues/10)) ([e1cb144](https://github.com/mistweaverco/kuba/commit/e1cb14464066ec6213bcff0673afd0b76256cf20))
* **docs:** update README with TOC ([7a5b37f](https://github.com/mistweaverco/kuba/commit/7a5b37fa49a3e2f78b21c435e82840e53d5c277e))
* **web:** overhaul website ([#9](https://github.com/mistweaverco/kuba/issues/9)) ([2d05d90](https://github.com/mistweaverco/kuba/commit/2d05d90ba3d2ea615b2f68f3d2e96ffa644d2afa))

# [0.5.0](https://github.com/mistweaverco/kuba/compare/v0.4.0...v0.5.0) (2025-08-25)


### Features

* **secret-path:** enforce all uppercase, add vale-config, add docs ([#8](https://github.com/mistweaverco/kuba/issues/8)) ([aeacc73](https://github.com/mistweaverco/kuba/commit/aeacc739e4812fa83dce48cc43bfbacb7e718c16))

# [0.4.0](https://github.com/mistweaverco/kuba/compare/v0.3.0...v0.4.0) (2025-08-24)


### Features

* **manager:** add secret-path ([#6](https://github.com/mistweaverco/kuba/issues/6)) ([594dd1a](https://github.com/mistweaverco/kuba/commit/594dd1a90b38ec92a6a128bc60d2cacb6e85527a))
* **managers:** add env interpolation ([#7](https://github.com/mistweaverco/kuba/issues/7)) ([4248f33](https://github.com/mistweaverco/kuba/commit/4248f337a9cb53ad62f7b3d5e1334de0d070bba8))

# [0.3.0](https://github.com/mistweaverco/kuba/compare/v0.2.0...v0.3.0) (2025-08-22)


### Bug Fixes

* **ci:** web deployment ([f4bd797](https://github.com/mistweaverco/kuba/commit/f4bd7978aaf083c1dc6d8a2a4136d2324b70886f))
* **docs:** fix typo ([a448d33](https://github.com/mistweaverco/kuba/commit/a448d33ec018fd88becf9815b6e061055504aab0))
* **web:** install script ([5cc8a33](https://github.com/mistweaverco/kuba/commit/5cc8a33fda80707691cf178e8042803640bea423))


### Features

* **managers:** add OpenBao support ([#5](https://github.com/mistweaverco/kuba/issues/5)) ([1177125](https://github.com/mistweaverco/kuba/commit/11771251a508495c798cbd672783c11a8e562e34))
* **web:** add website ([#4](https://github.com/mistweaverco/kuba/issues/4)) ([3ee56fc](https://github.com/mistweaverco/kuba/commit/3ee56fc743097977a37ed6c8c7940a50f53bf416))

# [0.2.0](https://github.com/mistweaverco/kuba/compare/v0.1.0...v0.2.0) (2025-08-20)


### Bug Fixes

* **docs:** README gcp id ([a6f0782](https://github.com/mistweaverco/kuba/commit/a6f078230a2130496fb01514402855e3cb77c139))

# [0.1.0](https://github.com/mistweaverco/kuba/compare/48ae181579d0981b1e5fc1e800b2c6f398e57306...v0.1.0) (2025-08-19)


### Bug Fixes

* **ci:** fix macos signing ([f52666c](https://github.com/mistweaverco/kuba/commit/f52666ca85e07703fafbe1c1fa3c8d6feca80b63))


### Features

* **app:** add bare app skeleton ([c921cfb](https://github.com/mistweaverco/kuba/commit/c921cfbe08bc2d5cd97ff4b879998c58a0c2a66a))
* **cmd:** add init subcommand ([b4c6df9](https://github.com/mistweaverco/kuba/commit/b4c6df908e7262ae199d5473a34129cf3f5e653f))
* **docs:** add schema ([6f931e5](https://github.com/mistweaverco/kuba/commit/6f931e54f2b94ec9c078e38e24a886a4e1f7f43c))
* **docs:** reworked README ([48ae181](https://github.com/mistweaverco/kuba/commit/48ae181579d0981b1e5fc1e800b2c6f398e57306))
* **Makefile:** add Makefile ([ee4a112](https://github.com/mistweaverco/kuba/commit/ee4a112b77800f0f2bbda9b84117c43dd4c4eec9))
* **secrets:** add gcp secret provider ([a428360](https://github.com/mistweaverco/kuba/commit/a428360b431335137a5387a70085dd104db4ac31))
