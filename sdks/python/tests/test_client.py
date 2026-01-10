"""Tests for renamed SDK client."""

import json
import pytest
import httpx
import respx

from renamed import (
    RenamedClient,
    AuthenticationError,
    ValidationError,
    RateLimitError,
    InsufficientCreditsError,
)
from renamed.client import AsyncJob


class TestRenamedClientInit:
    """Tests for client initialization."""

    def test_raises_without_api_key(self):
        """Should raise AuthenticationError when API key is missing."""
        with pytest.raises(AuthenticationError):
            RenamedClient(api_key="")

    def test_creates_with_valid_api_key(self):
        """Should create client with valid API key."""
        client = RenamedClient(api_key="rt_test123")
        assert client is not None

    def test_accepts_custom_options(self):
        """Should accept custom configuration options."""
        client = RenamedClient(
            api_key="rt_test123",
            base_url="https://custom.api.com",
            timeout=60.0,
            max_retries=5,
        )
        assert client._base_url == "https://custom.api.com"
        assert client._timeout == 60.0
        assert client._max_retries == 5


class TestRenamedClientErrors:
    """Tests for error handling."""

    @respx.mock
    def test_401_raises_authentication_error(self):
        """Should raise AuthenticationError on 401."""
        respx.get("https://www.renamed.to/api/v1/user").mock(
            return_value=httpx.Response(401, json={"error": "Invalid API key"})
        )

        client = RenamedClient(api_key="rt_invalid")

        with pytest.raises(AuthenticationError):
            client.get_user()

    @respx.mock
    def test_402_raises_insufficient_credits_error(self):
        """Should raise InsufficientCreditsError on 402."""
        respx.get("https://www.renamed.to/api/v1/user").mock(
            return_value=httpx.Response(402, json={"error": "Insufficient credits"})
        )

        client = RenamedClient(api_key="rt_test123")

        with pytest.raises(InsufficientCreditsError):
            client.get_user()

    @respx.mock
    def test_429_raises_rate_limit_error(self):
        """Should raise RateLimitError on 429."""
        respx.get("https://www.renamed.to/api/v1/user").mock(
            return_value=httpx.Response(
                429, json={"error": "Rate limit exceeded", "retryAfter": 60}
            )
        )

        client = RenamedClient(api_key="rt_test123")

        with pytest.raises(RateLimitError) as exc_info:
            client.get_user()

        assert exc_info.value.retry_after == 60

    @respx.mock
    def test_400_raises_validation_error(self):
        """Should raise ValidationError on 400."""
        respx.get("https://www.renamed.to/api/v1/user").mock(
            return_value=httpx.Response(400, json={"error": "Invalid request"})
        )

        client = RenamedClient(api_key="rt_test123")

        with pytest.raises(ValidationError):
            client.get_user()


class TestGetUser:
    """Tests for get_user method."""

    @respx.mock
    def test_returns_user_data(self):
        """Should return user data from API."""
        mock_user = {
            "id": "user123",
            "email": "test@example.com",
            "name": "Test User",
            "credits": 100,
        }

        respx.get("https://www.renamed.to/api/v1/user").mock(
            return_value=httpx.Response(200, json=mock_user)
        )

        client = RenamedClient(api_key="rt_test123")
        user = client.get_user()

        assert user.id == "user123"
        assert user.email == "test@example.com"
        assert user.name == "Test User"
        assert user.credits == 100


class TestRename:
    """Tests for rename method."""

    @respx.mock
    def test_renames_file_from_bytes(self):
        """Should rename file from bytes."""
        mock_result = {
            "originalFilename": "doc.pdf",
            "suggestedFilename": "2025-01-15_Invoice.pdf",
            "folderPath": "2025/Invoices",
            "confidence": 0.95,
        }

        respx.post("https://www.renamed.to/api/v1/rename").mock(
            return_value=httpx.Response(200, json=mock_result)
        )

        client = RenamedClient(api_key="rt_test123")
        result = client.rename(b"fake pdf content")

        assert result.original_filename == "doc.pdf"
        assert result.suggested_filename == "2025-01-15_Invoice.pdf"
        assert result.folder_path == "2025/Invoices"
        assert result.confidence == 0.95


class TestPdfSplit:
    """Tests for pdf_split method."""

    @respx.mock
    def test_returns_async_job(self):
        """Should return AsyncJob for polling."""
        respx.post("https://www.renamed.to/api/v1/pdf-split").mock(
            return_value=httpx.Response(
                200, json={"statusUrl": "https://api.example.com/status/job123"}
            )
        )

        client = RenamedClient(api_key="rt_test123")
        job = client.pdf_split(b"fake pdf content")

        assert job is not None
        assert hasattr(job, "wait")
        assert hasattr(job, "status")


class TestAsyncJob:
    """Tests for AsyncJob."""

    @respx.mock
    def test_polls_until_completed(self):
        """Should poll until job is completed."""
        mock_result = {
            "originalFilename": "multi.pdf",
            "documents": [
                {
                    "index": 0,
                    "filename": "doc1.pdf",
                    "pages": "1-5",
                    "downloadUrl": "https://...",
                    "size": 1000,
                }
            ],
            "totalPages": 10,
        }

        call_count = 0

        def status_callback(request):
            nonlocal call_count
            call_count += 1
            if call_count < 3:
                return httpx.Response(
                    200,
                    json={
                        "jobId": "job123",
                        "status": "processing",
                        "progress": call_count * 33,
                    },
                )
            return httpx.Response(
                200,
                json={
                    "jobId": "job123",
                    "status": "completed",
                    "progress": 100,
                    "result": mock_result,
                },
            )

        respx.get("https://api.example.com/status/job123").mock(side_effect=status_callback)

        client = RenamedClient(api_key="rt_test123")
        job = AsyncJob(client, "https://api.example.com/status/job123", 0.01, 10)

        progress_updates = []
        result = job.wait(lambda s: progress_updates.append(s.progress))

        assert result.original_filename == "multi.pdf"
        assert len(result.documents) == 1
        assert 33 in progress_updates
        assert 66 in progress_updates
