# Concierge

Concierge is a compliance check tool for AWS CloudFormation stack. The idea was previously implemented in [cfn-compliance-check](https://github.com/jc1518/cfn-compliance-check).

Why re-write it in Go?

- Learn some Go
- Write Once, Run Anywhere (WORA)

## Install

- Install from source: `go install github.com/jc1518/Concierge`
- Download compiled binary: [TODO]

## Usage

1. Setup CloudConformity API key environment variable `CLOUD_CONFORMITY_API_KEY` (You should be able to create one in CloudConformity console `User settings > API Keys` if you don't have one yet).

2. Setup your AWS credential (e.g. environment variables, profile or EC2 instance role).

3. Follow the usage. e.g `Concierge --stacks-arn arn:aws:cloudformation:ap-southeast-2:123456789000:stack/my-stack/69d48220-010d-11ec-982a-06dd10360dfc`

   ```
   Usage of Concierge:
   -stacks-arn string
           CloudFormation stacks ARN, use comma to seperate if more than one
   ```

## [TODO]

- Tests
- Github action to build and publish the binary
