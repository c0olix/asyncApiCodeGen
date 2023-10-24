## [2.1.1](https://github.com/c0olix/asyncApiCodeGen/compare/v2.1.0...v2.1.1) (2023-10-24)


### Bug Fixes

* missing import change on producer ([77c6646](https://github.com/c0olix/asyncApiCodeGen/commit/77c664643c5cf6f85042c4f4826cae0950e306fa))



# [2.1.0](https://github.com/c0olix/asyncApiCodeGen/compare/v2.0.2...v2.1.0) (2023-10-23)


### Bug Fixes

* missing go sum entries ([edd92e7](https://github.com/c0olix/asyncApiCodeGen/commit/edd92e75f768da8307f9463a7907832689aae967))
* package names ([ecd6b98](https://github.com/c0olix/asyncApiCodeGen/commit/ecd6b98502597b52419178fcf32e8738c1a2e88f))


### Features

* introduce springboot3 flavor for java ([968d9e9](https://github.com/c0olix/asyncApiCodeGen/commit/968d9e9a24db2551a819e84b4211acfaef236f5a))



## [2.0.2](https://github.com/c0olix/asyncApiCodeGen/compare/v2.0.1...v2.0.2) (2023-03-24)


### Bug Fixes

* import time when format is date ([2f2dab8](https://github.com/c0olix/asyncApiCodeGen/commit/2f2dab80d19df479fc99c32fb6c2180cfef87983))



## [2.0.1](https://github.com/c0olix/asyncApiCodeGen/compare/v2.0.0...v2.0.1) (2023-02-02)


### Bug Fixes

* only import context if has producer ([ce3a29c](https://github.com/c0olix/asyncApiCodeGen/commit/ce3a29c3330d4e585ffa6c8031c5214c13f564d1))



# [2.0.0](https://github.com/c0olix/asyncApiCodeGen/compare/v1.2.0...v2.0.0) (2023-01-31)


* feat!: new version of goChan (including breaking changes) ([fa0e046](https://github.com/c0olix/asyncApiCodeGen/commit/fa0e0468b75801d2128b43bd23f4dea6c7f304d0))


### BREAKING CHANGES

* generated code will break old versions



# [1.2.0](https://github.com/c0olix/asyncApiCodeGen/compare/v1.1.0...v1.2.0) (2023-01-22)


### Features

* info logging ([f4f3d16](https://github.com/c0olix/asyncApiCodeGen/commit/f4f3d16694af3de822338c9985201adcb2535c39))



# [1.1.0](https://github.com/c0olix/asyncApiCodeGen/compare/v1.0.3...v1.1.0) (2023-01-22)


### Features

* add MQTT as flavor ([58cfe27](https://github.com/c0olix/asyncApiCodeGen/commit/58cfe274f03881764335526ccf894c742197d309))



## [1.0.3](https://github.com/c0olix/asyncApiCodeGen/compare/v1.0.2...v1.0.3) (2023-01-12)


### Bug Fixes

* add missing import ([e13d6c4](https://github.com/c0olix/asyncApiCodeGen/commit/e13d6c4a2c86b6e74a7684b79124940267a674ff))



## [1.0.2](https://github.com/c0olix/asyncApiCodeGen/compare/v1.0.1...v1.0.2) (2023-01-04)


### Bug Fixes

* add missing import ([ce3a8aa](https://github.com/c0olix/asyncApiCodeGen/commit/ce3a8aa224b703889ade56a8ffa332e58ba36d03))



## [1.0.1](https://github.com/c0olix/asyncApiCodeGen/compare/v1.0.0...v1.0.1) (2023-01-03)


### Bug Fixes

* remove unused imports from template ([868e58d](https://github.com/c0olix/asyncApiCodeGen/commit/868e58d8ca13d048b95627faebcbfa79417943ab))



# [1.0.0](https://github.com/c0olix/asyncApiCodeGen/compare/v0.7.0...v1.0.0) (2022-09-04)


* feat(cli)!: add flag packageName ([9186a5d](https://github.com/c0olix/asyncApiCodeGen/commit/9186a5dcde97ceff911e79c0d6803e8f85a21bec))
* feat(cli)!: add flags input, output ([24c473d](https://github.com/c0olix/asyncApiCodeGen/commit/24c473d3a9bb5daa7cd13766fc9d3636be8edb05))
* feat(cli)!: add flag createDir ([b1574fb](https://github.com/c0olix/asyncApiCodeGen/commit/b1574fb5cf603634ade6840238b1e4f3be41488d))


### Features

* **go:** create a not private version for go ([cdb6373](https://github.com/c0olix/asyncApiCodeGen/commit/cdb63739f92674461ec3fc0dcb30a10930cfb2f3))


### BREAKING CHANGES

* without given packageName the cli fails
* input and output location via arguments aren't supported anymore
* java codegen does not create output dir by default anymore



# [0.7.0](https://github.com/c0olix/asyncApiCodeGen/compare/v0.6.2...v0.7.0) (2022-09-04)


### Bug Fixes

* remove example output, as it uses a private repository and breaks tests ([44b6019](https://github.com/c0olix/asyncApiCodeGen/commit/44b601939c90d16633a46d3cda5d6816b55f2173))
* remove password type, as it isn't supported ([e4e5f9d](https://github.com/c0olix/asyncApiCodeGen/commit/e4e5f9d25276ec12b75824e0a60a2785625841c2))
* remove private repository dependency ([79b9baa](https://github.com/c0olix/asyncApiCodeGen/commit/79b9baa73789d27fb008f4b7439b6699c816f790))


### Features

* **cli:** validate given async api spec ([4352cbf](https://github.com/c0olix/asyncApiCodeGen/commit/4352cbf5a85fe2e956032b096ba77d763a23a073))
* **mosaicGoKafkaGenerator:** more validations ([0e8ae54](https://github.com/c0olix/asyncApiCodeGen/commit/0e8ae54858d2c2eb05fb3983aa163423f3f5cd6d))
* **mosaicJavaKafkaGenerator:** more validations ([01b727b](https://github.com/c0olix/asyncApiCodeGen/commit/01b727b0c03b35a624083396a2b41745d00e9fcd))



## [0.6.2](https://github.com/c0olix/asyncApiCodeGen/compare/v0.6.1...v0.6.2) (2022-08-31)


### Bug Fixes

* don't validate required on bool fields ([43ec5c9](https://github.com/c0olix/asyncApiCodeGen/commit/43ec5c9907cb266f7a3a474532bbc6ae0ce8ce70))



## [0.6.1](https://github.com/c0olix/asyncApiCodeGen/compare/v0.5.0...v0.6.1) (2022-08-31)


### Bug Fixes

* **mosaicGoKafkaGenerator:** export methods ([b02124b](https://github.com/c0olix/asyncApiCodeGen/commit/b02124bafb620385eb34a6425649b316f91f984f))
* **mosaicJavaKafkaGenerator:** create output dir if not present ([2060ea4](https://github.com/c0olix/asyncApiCodeGen/commit/2060ea4e01bb73079800ee31780e104f08264250))
* **mosaicJavaKafkaGenerator:** imports for array items and generation of nested objects ([0d87f6b](https://github.com/c0olix/asyncApiCodeGen/commit/0d87f6be6dcece40d91b672ed3e323a16f5d39dc))
* naming and pointer ([6b64c55](https://github.com/c0olix/asyncApiCodeGen/commit/6b64c55f2fc3481d6f44eb64a132db1867ba6674))


### Features

* **mosaicGoKafkaGenerator:** validation on struct fields ([cc51860](https://github.com/c0olix/asyncApiCodeGen/commit/cc518606b6d7f252e49202b5220b61f8314374c8))



# [0.5.0](https://github.com/c0olix/asyncApiCodeGen/compare/v0.4.0...v0.5.0) (2022-08-29)


### Bug Fixes

* don't inline objects and items ([8ef53cf](https://github.com/c0olix/asyncApiCodeGen/commit/8ef53cf393478435cab09dc8fee207ebbe108daa))
* time import for go ([548da5f](https://github.com/c0olix/asyncApiCodeGen/commit/548da5f9e0ec15c93faba2bf42e4eff81cd90458))


### Features

* add default value ([faf334e](https://github.com/c0olix/asyncApiCodeGen/commit/faf334e8e17a2f7faa8fa8cb7749f3c98be4d409))
* use async api parser dependency ([371f028](https://github.com/c0olix/asyncApiCodeGen/commit/371f0286a611d80ba0f95976001bbb09103a37e5))



# [0.4.0](https://github.com/c0olix/asyncApiCodeGen/compare/v0.3.3...v0.4.0) (2022-08-23)


### Features

* import for java ([e186c97](https://github.com/c0olix/asyncApiCodeGen/commit/e186c97f53ef1302ce1f9881e7229f3302fb246f))
* respect opId for java ([6c01f53](https://github.com/c0olix/asyncApiCodeGen/commit/6c01f53d72dc7f3d58c5ff43db8f4395e98d9813))



## [0.3.3](https://github.com/c0olix/asyncApiCodeGen/compare/v0.3.0...v0.3.3) (2022-08-23)


### Bug Fixes

* array object conversation ([fecea88](https://github.com/c0olix/asyncApiCodeGen/commit/fecea8816f1a6cbb31da4857754b242df7e9a417))
* respect OperationId ([01f2039](https://github.com/c0olix/asyncApiCodeGen/commit/01f203939997a9d7bcf0f41d9385a9c7fa564df8))



# [0.3.0](https://github.com/c0olix/asyncApiCodeGen/compare/v0.2.0...v0.3.0) (2022-08-22)


### Features

* java items can also have ref string ([564db1b](https://github.com/c0olix/asyncApiCodeGen/commit/564db1bdb64a201afd2b48877d546b4470fd50a6))



# [0.2.0](https://github.com/c0olix/asyncApiCodeGen/compare/v0.1.3...v0.2.0) (2022-08-22)


### Features

* array items may have a reference string ([a621477](https://github.com/c0olix/asyncApiCodeGen/commit/a62147791b1efe12a9599bafa245946d926086f3))



## [0.1.3](https://github.com/c0olix/asyncApiCodeGen/compare/78e9d2c1e1252aaf0a6ab0cbd81bca4b309791b4...v0.1.3) (2022-08-22)


### Bug Fixes

* correct format ([f48ce91](https://github.com/c0olix/asyncApiCodeGen/commit/f48ce912a6631ca80c93010717919f9ad9e418d0))
* go version ([b711a7a](https://github.com/c0olix/asyncApiCodeGen/commit/b711a7a3f5b379c79f6dca61c12639872a8db995))


### Features

* consumer interface for java ([436d130](https://github.com/c0olix/asyncApiCodeGen/commit/436d130f6885de74b194efc65e9bcbbec2225cca))
* create documentation commentaries ([3683694](https://github.com/c0olix/asyncApiCodeGen/commit/368369466da04102465bb51af38a3181a0365952))
* embed template files and some minor fixes ([207cba0](https://github.com/c0olix/asyncApiCodeGen/commit/207cba0b4c8c7a965ed199204c5d2e0a6a16d7e7))
* format output ([3ddf78f](https://github.com/c0olix/asyncApiCodeGen/commit/3ddf78ff71580070324701388ffcc9f9752ea7e7))
* generate command ([a09bdb2](https://github.com/c0olix/asyncApiCodeGen/commit/a09bdb28f7dd7a46e58db5c10cf74edbfbb53037))
* java code generator (types) ([4e43323](https://github.com/c0olix/asyncApiCodeGen/commit/4e4332323226d84c9198f558711e6067df5470e9))
* logging ([0e4c7a6](https://github.com/c0olix/asyncApiCodeGen/commit/0e4c7a67a54a7926fdb83e6f45295c2051fc25d2))
* producer interface for java ([07742d7](https://github.com/c0olix/asyncApiCodeGen/commit/07742d78b8028829b9301e67d5d1d7bb622e1f1a))
* root command ([78e9d2c](https://github.com/c0olix/asyncApiCodeGen/commit/78e9d2c1e1252aaf0a6ab0cbd81bca4b309791b4))
* separate different generators ([2123c45](https://github.com/c0olix/asyncApiCodeGen/commit/2123c4581488dd92257f36d87b1a9e77afb367ef))
* type for asyncApSpec with conversion to go specific syntax ([e3154f6](https://github.com/c0olix/asyncApiCodeGen/commit/e3154f60a7eff2a6d0a80a78730a5668dc25d95f))



