"""Functionality for interacting with the subscriptions API."""
from typing import AsyncIterator, Optional, Set

from ..constants import PLANET_BASE_URL
from ..http import Session

BASE_URL = f'{PLANET_BASE_URL}/subscriptions/v1'


class SubscriptionsClient:
    """Low-level asynchronous access to Planet's subscriptions API."""

    def __init__(self, session: Session, base_url: str = None):
        """Initialize the client.

        Parameters:
            session: Open session connected to server.
            base_url: The base URL to use. Defaults to production subscriptions
                API base URL.
        """
        self._session = session

        self._base_url = base_url or BASE_URL
        if self._base_url.endswith('/'):
            self._base_url = self._base_url[:-1]

    async def list_subscriptions(
            self,
            status: Optional[Set[str]] = None,
            limit: int = 100) -> AsyncIterator[dict]:
        """List subscriptions.

        Parameters:
            status: Filter subscriptions by status.
            limit: Maximum number of subscriptions to return.

        Yields:
            Description of a subscription.
        """
        raise NotImplementedError

    async def create_subscription(self, request: dict) -> dict:
        """Create a subscription.

        Parameters:
            request: Subscription request definition.

        Returns:
            Description of the created subscription.
        """
        raise NotImplementedError

    async def cancel_subscription(self, subscription_id: str) -> dict:
        """Cancel a subscription.

        Parameters:
            subscription_id: Subscription identifier.

        Returns:
            Description of the cancelled subscription.
        """
        raise NotImplementedError

    async def update_subscription(self, subscription_id: str,
                                  request: dict) -> dict:
        """Update a subscription.

        Parameters:
            subscription_id: Subscription identifier.
            request: Subscription request definition.

        Returns:
            Description of the updated subscription.
        """
        raise NotImplementedError

    async def get_subscription(self, subscription_id: str) -> dict:
        """Get a subscription.

        Parameters:
            subscription_id: Subscription identifier.

        Returns:
            Description of the subscription.
        """
        raise NotImplementedError

    async def get_results(
            self,
            subscription_id: str,
            status: Optional[Set[str]] = None,
            limit: int = 100) -> AsyncIterator[dict]:
        """List results for a subscription.

        Parameters:
            subscription_id: Subscription identifier.
            status: Filter results by status.
            limit: Maximum number of results to return.

        Yields:
            Description of a subscription result.
        """
        raise NotImplementedError
