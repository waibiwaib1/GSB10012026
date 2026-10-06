"""
Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
SPDX-License-Identifier: MIT-0
"""

from cfnlint.rules import CloudFormationLintRule, RuleMatch


class CacheClusterEngineVersion(CloudFormationLintRule):
    """Check ElastiCache Redis engine version format"""

    id = "E3059"
    shortdesc = "Validate ElastiCache Redis engine version"
    description = (
        "ElastiCache Redis versions 6 and later must be specified using the"
        " major and minor version"
    )
    source_url = (
        "https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/"
        "aws-properties-elasticache-cache-cluster.html#cfn-elasticache-cachecluster-engineversion"
    )
    tags = ["resources", "elasticcache"]

    def __init__(self):
        """Init"""
        super().__init__()
        self.resource_property_types.append("AWS::ElastiCache::CacheCluster")

    def _match_properties(self, properties, path, scenario=None):
        matches = []
        engine = properties.get("Engine")
        engine_version = properties.get("EngineVersion")

        if not isinstance(engine, str) or not isinstance(engine_version, str):
            return matches

        if engine.lower() != "redis":
            return matches

        version_parts = engine_version.split(".")
        if len(version_parts) < 3:
            return matches

        try:
            major_version = int(version_parts[0])
        except ValueError:
            return matches

        if major_version < 6:
            return matches

        pathmessage = list(path) + ["EngineVersion"]
        if scenario is not None:
            scenario_text = " and ".join(
                [
                    f'when condition "{key}" is {value}'
                    for key, value in scenario.items()
                ]
            )
            message = (
                '"EngineVersion" for Redis version 6 and later must not include a'
                " patch version when {0} at {1}"
            )
            message_args = (scenario_text, "/".join(map(str, pathmessage)))
        else:
            message = (
                '"EngineVersion" for Redis version 6 and later must not include a'
                " patch version at {0}"
            )
            message_args = ("/".join(map(str, pathmessage)),)

        matches.append(RuleMatch(pathmessage, message.format(*message_args)))
        return matches

    def match_resource_properties(self, properties, _, path, cfn):
        """Check CloudFormation Properties"""
        matches = []
        scenarios = cfn.get_conditions_scenarios_from_object([properties])
        if scenarios:
            for scenario in scenarios:
                matches.extend(
                    self._match_properties(
                        cfn.get_value_from_scenario(properties, scenario),
                        path,
                        scenario,
                    )
                )
        else:
            matches.extend(self._match_properties(properties, path))

        return matches
