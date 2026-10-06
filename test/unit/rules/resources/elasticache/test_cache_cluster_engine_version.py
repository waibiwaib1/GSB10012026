"""
Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
SPDX-License-Identifier: MIT-0
"""

from collections import deque

from cfnlint import Template
from cfnlint.rules.resources.elasticache.CacheClusterEngineVersion import (
    CacheClusterEngineVersion,
)
from test.testlib.testcase import BaseTestCase


class TestCacheClusterEngineVersion(BaseTestCase):
    """Test ElastiCache CacheCluster EngineVersion"""

    def setUp(self):
        super().setUp()
        self.rule = CacheClusterEngineVersion()
        self.path = deque(["Resources", "RedisCluster", "Properties"])
        self.cfn = Template("", {})

    def test_invalid_redis_patch_version(self):
        matches = self.rule.match_resource_properties(
            {"Engine": "redis", "EngineVersion": "7.1.0"},
            "AWS::ElastiCache::CacheCluster",
            self.path,
            self.cfn,
        )

        self.assertEqual(1, len(matches))
        self.assertEqual(
            ["Resources", "RedisCluster", "Properties", "EngineVersion"],
            matches[0].path,
        )

    def test_valid_redis_minor_version(self):
        matches = self.rule.match_resource_properties(
            {"Engine": "redis", "EngineVersion": "7.1"},
            "AWS::ElastiCache::CacheCluster",
            self.path,
            self.cfn,
        )

        self.assertEqual([], matches)

    def test_legacy_redis_patch_version(self):
        matches = self.rule.match_resource_properties(
            {"Engine": "redis", "EngineVersion": "5.0.6"},
            "AWS::ElastiCache::CacheCluster",
            self.path,
            self.cfn,
        )

        self.assertEqual([], matches)

    def test_non_redis_patch_version(self):
        matches = self.rule.match_resource_properties(
            {"Engine": "memcached", "EngineVersion": "1.6.17"},
            "AWS::ElastiCache::CacheCluster",
            self.path,
            self.cfn,
        )

        self.assertEqual([], matches)

    def test_unresolved_engine_version(self):
        matches = self.rule.match_resource_properties(
            {"Engine": "redis", "EngineVersion": {"Ref": "EngineVersion"}},
            "AWS::ElastiCache::CacheCluster",
            self.path,
            self.cfn,
        )

        self.assertEqual([], matches)
