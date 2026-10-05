"""Functionality for interacting with the subscriptions API."""
import typing

from ..constants import PLANET_BASE_URL
from ..http import Session

BASE_URL = f'{PLANET_BASE_URL}/subscriptions/v1'


class SubscriptionsClient:
    """High-level asynchronous access to Planet's subscriptions API."""

    def __init__(self, session: Session, base_url: str = None):
        """Initialize the client.

        Parameters:
            session: Open session connected to server.
            base_url: The base URL to use. Defaults to production
                subscriptions API base url.
        """
        self._session = session
        self._base_url = base_url or BASE_URL

        if self._base_url.endswith('/'):
            self._base_url = self._base_url[:-1]

    async def create_subscription(self, request: dict) -> dict:
        """Create a subscription.

        Parameters:
            request: Subscription request definition.

        Returns:
            JSON description of the created subscription.
        """
        raise NotImplementedError

    async def update_subscription(self, subscription_id: str,
                                  request: dict) -> dict:
        """Update a subscription.

        Parameters:
            subscription_id: The ID of the subscription.
            request: Subscription request definition.

        Returns:
            JSON description of the updated subscription.
        """
        raise NotImplementedError

    async def get_results(
            self,
            subscription_id: str,
            status: typing.Set[str] = None,
            limit: typing.Union[int, None] = 100
    ) -> typing.AsyncIterator[dict]:
        """Get results for a subscription.

        Parameters:
            subscription_id: The ID of the subscription.
            status: Filter results by status.
            limit: Maximum number of results to return.

        Returns:
            An iterator over subscription results.
        """
        raise NotImplementedError  # pragma: no cover
        yield  # pragma: no cover

    async def get_subscription(self, subscription_id: str) -> dict:
        """Get subscription details.

        Parameters:
            subscription_id: The ID of the subscription.

        Returns:
            JSON description of the subscription.
        """
        raise NotImplementedError

    async def cancel_subscription(self, subscription_id: str) -> dict:
        """Cancel a subscription.

        Parameters:
            subscription_id: The ID of the subscription.

        Returns:
            Results of the cancel request.
        """
        raise NotImplementedError

    async def list_subscriptions(
            self,
            status: typing.Set[str] = None,
            limit: typing.Union[int, None] = 100
    ) -> typing.AsyncIterator[dict]:
        """List subscriptions.

        Parameters:
            status: Filter subscriptions by status.
            limit: Maximum number of subscriptions to return.

        Returns:
            An iterator over subscriptions.
        """
        raise NotImplementedError  # pragma: no cover
        yield  # pragma: no cover
