# SIPL Conformance Suite

This directory contains an independent test suite for verifying the correctness of SIPL language implementations.

## `suite.json` Format
The tests are described as an array of JSON objects. Each object has the following structure:
- `name` (string): The name of the test.
- `source` (string): The source code of the SIPL program.
- `expected_state` (map[string]int): The expected values of the variables after execution (ignored if an error is expected).
- `expect_error` (boolean): `true` if the program contains a syntax or runtime error.

Translated with DeepL.com (free version)