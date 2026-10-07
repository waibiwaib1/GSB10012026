<?php

declare(strict_types=1);

namespace Tests\Psalm\LaravelPlugin\Unit\Handlers\Auth;

use PhpParser\Node\Expr\MethodCall;
use PhpParser\NodeFinder;
use PHPUnit\Framework\TestCase;
use Psalm\Internal\Provider\StatementsProvider;
use Psalm\LaravelPlugin\Handlers\Auth\GuardHandler;

use const PHP_VERSION_ID;

/** @covers \Psalm\LaravelPlugin\Handlers\Auth\GuardHandler */
final class GuardHandlerTest extends TestCase
{
    /** @test */
    public function it_finds_an_explicit_guard_after_an_auth_helper_call(): void
    {
        $methodCall = $this->findMethodCall(<<<'PHP'
<?php

$user = auth('admin')->user();
PHP);

        $findGuardName = new \ReflectionMethod(GuardHandler::class, 'findGuardNameInCallChain');
        $findGuardName->setAccessible(true);

        self::assertSame('admin', $findGuardName->invoke(null, $methodCall));
    }

    /** @test */
    public function it_matches_the_auth_helper_case_insensitively(): void
    {
        $methodCall = $this->findMethodCall(<<<'PHP'
<?php

$user = AUTH('admin')->user();
PHP);

        $findGuardName = new \ReflectionMethod(GuardHandler::class, 'findGuardNameInCallChain');
        $findGuardName->setAccessible(true);

        self::assertSame('admin', $findGuardName->invoke(null, $methodCall));
    }

    private function findMethodCall(string $source): MethodCall
    {
        $hasErrors = false;
        $statements = StatementsProvider::parseStatements($source, PHP_VERSION_ID, $hasErrors);

        self::assertFalse($hasErrors);

        $methodCalls = (new NodeFinder())->findInstanceOf($statements, MethodCall::class);

        self::assertCount(1, $methodCalls);
        self::assertInstanceOf(MethodCall::class, $methodCalls[0]);

        return $methodCalls[0];
    }
}
