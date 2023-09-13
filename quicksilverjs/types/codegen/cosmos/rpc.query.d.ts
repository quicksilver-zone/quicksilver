export declare const createRPCQueryClient: ({ rpcEndpoint }: {
    rpcEndpoint: string;
}) => Promise<{
    cosmos: {
        auth: {
            v1beta1: {
                accounts(request?: import("./auth/v1beta1/query").QueryAccountsRequest): Promise<import("./auth/v1beta1/query").QueryAccountsResponse>;
                account(request: import("./auth/v1beta1/query").QueryAccountRequest): Promise<import("./auth/v1beta1/query").QueryAccountResponse>;
                params(request?: import("./auth/v1beta1/query").QueryParamsRequest): Promise<import("./auth/v1beta1/query").QueryParamsResponse>;
            };
        };
        authz: {
            v1beta1: {
                grants(request: import("./authz/v1beta1/query").QueryGrantsRequest): Promise<import("./authz/v1beta1/query").QueryGrantsResponse>;
            };
        };
        bank: {
            v1beta1: {
                balance(request: import("./bank/v1beta1/query").QueryBalanceRequest): Promise<import("./bank/v1beta1/query").QueryBalanceResponse>;
                allBalances(request: import("./bank/v1beta1/query").QueryAllBalancesRequest): Promise<import("./bank/v1beta1/query").QueryAllBalancesResponse>;
                totalSupply(request?: import("./bank/v1beta1/query").QueryTotalSupplyRequest): Promise<import("./bank/v1beta1/query").QueryTotalSupplyResponse>;
                supplyOf(request: import("./bank/v1beta1/query").QuerySupplyOfRequest): Promise<import("./bank/v1beta1/query").QuerySupplyOfResponse>;
                params(request?: import("./bank/v1beta1/query").QueryParamsRequest): Promise<import("./bank/v1beta1/query").QueryParamsResponse>;
                denomMetadata(request: import("./bank/v1beta1/query").QueryDenomMetadataRequest): Promise<import("./bank/v1beta1/query").QueryDenomMetadataResponse>;
                denomsMetadata(request?: import("./bank/v1beta1/query").QueryDenomsMetadataRequest): Promise<import("./bank/v1beta1/query").QueryDenomsMetadataResponse>;
            };
        };
        base: {
            tendermint: {
                v1beta1: {
                    getNodeInfo(request?: import("./base/tendermint/v1beta1/query").GetNodeInfoRequest): Promise<import("./base/tendermint/v1beta1/query").GetNodeInfoResponse>;
                    getSyncing(request?: import("./base/tendermint/v1beta1/query").GetSyncingRequest): Promise<import("./base/tendermint/v1beta1/query").GetSyncingResponse>;
                    getLatestBlock(request?: import("./base/tendermint/v1beta1/query").GetLatestBlockRequest): Promise<import("./base/tendermint/v1beta1/query").GetLatestBlockResponse>;
                    getBlockByHeight(request: import("./base/tendermint/v1beta1/query").GetBlockByHeightRequest): Promise<import("./base/tendermint/v1beta1/query").GetBlockByHeightResponse>;
                    getLatestValidatorSet(request?: import("./base/tendermint/v1beta1/query").GetLatestValidatorSetRequest): Promise<import("./base/tendermint/v1beta1/query").GetLatestValidatorSetResponse>;
                    getValidatorSetByHeight(request: import("./base/tendermint/v1beta1/query").GetValidatorSetByHeightRequest): Promise<import("./base/tendermint/v1beta1/query").GetValidatorSetByHeightResponse>;
                };
            };
        };
        distribution: {
            v1beta1: {
                params(request?: import("./distribution/v1beta1/query").QueryParamsRequest): Promise<import("./distribution/v1beta1/query").QueryParamsResponse>;
                validatorOutstandingRewards(request: import("./distribution/v1beta1/query").QueryValidatorOutstandingRewardsRequest): Promise<import("./distribution/v1beta1/query").QueryValidatorOutstandingRewardsResponse>;
                validatorCommission(request: import("./distribution/v1beta1/query").QueryValidatorCommissionRequest): Promise<import("./distribution/v1beta1/query").QueryValidatorCommissionResponse>;
                validatorSlashes(request: import("./distribution/v1beta1/query").QueryValidatorSlashesRequest): Promise<import("./distribution/v1beta1/query").QueryValidatorSlashesResponse>;
                delegationRewards(request: import("./distribution/v1beta1/query").QueryDelegationRewardsRequest): Promise<import("./distribution/v1beta1/query").QueryDelegationRewardsResponse>;
                delegationTotalRewards(request: import("./distribution/v1beta1/query").QueryDelegationTotalRewardsRequest): Promise<import("./distribution/v1beta1/query").QueryDelegationTotalRewardsResponse>;
                delegatorValidators(request: import("./distribution/v1beta1/query").QueryDelegatorValidatorsRequest): Promise<import("./distribution/v1beta1/query").QueryDelegatorValidatorsResponse>;
                delegatorWithdrawAddress(request: import("./distribution/v1beta1/query").QueryDelegatorWithdrawAddressRequest): Promise<import("./distribution/v1beta1/query").QueryDelegatorWithdrawAddressResponse>;
                communityPool(request?: import("./distribution/v1beta1/query").QueryCommunityPoolRequest): Promise<import("./distribution/v1beta1/query").QueryCommunityPoolResponse>;
            };
        };
        evidence: {
            v1beta1: {
                evidence(request: import("./evidence/v1beta1/query").QueryEvidenceRequest): Promise<import("./evidence/v1beta1/query").QueryEvidenceResponse>;
                allEvidence(request?: import("./evidence/v1beta1/query").QueryAllEvidenceRequest): Promise<import("./evidence/v1beta1/query").QueryAllEvidenceResponse>;
            };
        };
        feegrant: {
            v1beta1: {
                allowance(request: import("./feegrant/v1beta1/query").QueryAllowanceRequest): Promise<import("./feegrant/v1beta1/query").QueryAllowanceResponse>;
                allowances(request: import("./feegrant/v1beta1/query").QueryAllowancesRequest): Promise<import("./feegrant/v1beta1/query").QueryAllowancesResponse>;
            };
        };
        gov: {
            v1beta1: {
                proposal(request: import("./gov/v1beta1/query").QueryProposalRequest): Promise<import("./gov/v1beta1/query").QueryProposalResponse>;
                proposals(request: import("./gov/v1beta1/query").QueryProposalsRequest): Promise<import("./gov/v1beta1/query").QueryProposalsResponse>;
                vote(request: import("./gov/v1beta1/query").QueryVoteRequest): Promise<import("./gov/v1beta1/query").QueryVoteResponse>;
                votes(request: import("./gov/v1beta1/query").QueryVotesRequest): Promise<import("./gov/v1beta1/query").QueryVotesResponse>;
                params(request: import("./gov/v1beta1/query").QueryParamsRequest): Promise<import("./gov/v1beta1/query").QueryParamsResponse>;
                deposit(request: import("./gov/v1beta1/query").QueryDepositRequest): Promise<import("./gov/v1beta1/query").QueryDepositResponse>;
                deposits(request: import("./gov/v1beta1/query").QueryDepositsRequest): Promise<import("./gov/v1beta1/query").QueryDepositsResponse>;
                tallyResult(request: import("./gov/v1beta1/query").QueryTallyResultRequest): Promise<import("./gov/v1beta1/query").QueryTallyResultResponse>;
            };
        };
        mint: {
            v1beta1: {
                params(request?: import("./mint/v1beta1/query").QueryParamsRequest): Promise<import("./mint/v1beta1/query").QueryParamsResponse>;
                inflation(request?: import("./mint/v1beta1/query").QueryInflationRequest): Promise<import("./mint/v1beta1/query").QueryInflationResponse>;
                annualProvisions(request?: import("./mint/v1beta1/query").QueryAnnualProvisionsRequest): Promise<import("./mint/v1beta1/query").QueryAnnualProvisionsResponse>;
            };
        };
        params: {
            v1beta1: {
                params(request: import("./params/v1beta1/query").QueryParamsRequest): Promise<import("./params/v1beta1/query").QueryParamsResponse>;
            };
        };
        slashing: {
            v1beta1: {
                params(request?: import("./slashing/v1beta1/query").QueryParamsRequest): Promise<import("./slashing/v1beta1/query").QueryParamsResponse>;
                signingInfo(request: import("./slashing/v1beta1/query").QuerySigningInfoRequest): Promise<import("./slashing/v1beta1/query").QuerySigningInfoResponse>;
                signingInfos(request?: import("./slashing/v1beta1/query").QuerySigningInfosRequest): Promise<import("./slashing/v1beta1/query").QuerySigningInfosResponse>;
            };
        };
        staking: {
            v1beta1: {
                validators(request: import("./staking/v1beta1/query").QueryValidatorsRequest): Promise<import("./staking/v1beta1/query").QueryValidatorsResponse>;
                validator(request: import("./staking/v1beta1/query").QueryValidatorRequest): Promise<import("./staking/v1beta1/query").QueryValidatorResponse>;
                validatorDelegations(request: import("./staking/v1beta1/query").QueryValidatorDelegationsRequest): Promise<import("./staking/v1beta1/query").QueryValidatorDelegationsResponse>;
                validatorUnbondingDelegations(request: import("./staking/v1beta1/query").QueryValidatorUnbondingDelegationsRequest): Promise<import("./staking/v1beta1/query").QueryValidatorUnbondingDelegationsResponse>;
                delegation(request: import("./staking/v1beta1/query").QueryDelegationRequest): Promise<import("./staking/v1beta1/query").QueryDelegationResponse>;
                unbondingDelegation(request: import("./staking/v1beta1/query").QueryUnbondingDelegationRequest): Promise<import("./staking/v1beta1/query").QueryUnbondingDelegationResponse>;
                delegatorDelegations(request: import("./staking/v1beta1/query").QueryDelegatorDelegationsRequest): Promise<import("./staking/v1beta1/query").QueryDelegatorDelegationsResponse>;
                delegatorUnbondingDelegations(request: import("./staking/v1beta1/query").QueryDelegatorUnbondingDelegationsRequest): Promise<import("./staking/v1beta1/query").QueryDelegatorUnbondingDelegationsResponse>;
                redelegations(request: import("./staking/v1beta1/query").QueryRedelegationsRequest): Promise<import("./staking/v1beta1/query").QueryRedelegationsResponse>;
                delegatorValidators(request: import("./staking/v1beta1/query").QueryDelegatorValidatorsRequest): Promise<import("./staking/v1beta1/query").QueryDelegatorValidatorsResponse>;
                delegatorValidator(request: import("./staking/v1beta1/query").QueryDelegatorValidatorRequest): Promise<import("./staking/v1beta1/query").QueryDelegatorValidatorResponse>;
                historicalInfo(request: import("./staking/v1beta1/query").QueryHistoricalInfoRequest): Promise<import("./staking/v1beta1/query").QueryHistoricalInfoResponse>;
                pool(request?: import("./staking/v1beta1/query").QueryPoolRequest): Promise<import("./staking/v1beta1/query").QueryPoolResponse>;
                params(request?: import("./staking/v1beta1/query").QueryParamsRequest): Promise<import("./staking/v1beta1/query").QueryParamsResponse>;
            };
        };
        tx: {
            v1beta1: {
                simulate(request: import("./tx/v1beta1/service").SimulateRequest): Promise<import("./tx/v1beta1/service").SimulateResponse>;
                getTx(request: import("./tx/v1beta1/service").GetTxRequest): Promise<import("./tx/v1beta1/service").GetTxResponse>;
                broadcastTx(request: import("./tx/v1beta1/service").BroadcastTxRequest): Promise<import("./tx/v1beta1/service").BroadcastTxResponse>;
                getTxsEvent(request: import("./tx/v1beta1/service").GetTxsEventRequest): Promise<import("./tx/v1beta1/service").GetTxsEventResponse>;
            };
        };
        upgrade: {
            v1beta1: {
                currentPlan(request?: import("./upgrade/v1beta1/query").QueryCurrentPlanRequest): Promise<import("./upgrade/v1beta1/query").QueryCurrentPlanResponse>;
                appliedPlan(request: import("./upgrade/v1beta1/query").QueryAppliedPlanRequest): Promise<import("./upgrade/v1beta1/query").QueryAppliedPlanResponse>;
                upgradedConsensusState(request: import("./upgrade/v1beta1/query").QueryUpgradedConsensusStateRequest): Promise<import("./upgrade/v1beta1/query").QueryUpgradedConsensusStateResponse>;
                moduleVersions(request: import("./upgrade/v1beta1/query").QueryModuleVersionsRequest): Promise<import("./upgrade/v1beta1/query").QueryModuleVersionsResponse>;
            };
        };
    };
}>;
