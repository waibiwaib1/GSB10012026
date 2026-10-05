<?php

/*
 * This file is part of PHP CS Fixer.
 *
 * (c) Fabien Potencier <fabien@symfony.com>
 *     Dariusz Rumiński <dariusz.ruminski@gmail.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace PhpCsFixer\Tests\Console\Command;

use PhpCsFixer\Console\Application;
use PhpCsFixer\Console\Command\ListSetsCommand;
use PhpCsFixer\RuleSet\RuleSets;
use PhpCsFixer\Tests\TestCase;
use Symfony\Component\Console\Tester\CommandTester;

/**
 * @internal
 *
 * @covers \PhpCsFixer\Console\Command\ListSetsCommand
 */
final class ListSetsCommandTest extends TestCase
{
    public function testListSets()
    {
        $application = new Application();
        $application->add(new ListSetsCommand());

        $command = $application->find('list-sets');
        $commandTester = new CommandTester($command);

        $commandTester->execute([]);

        static::assertSame(0, $commandTester->getStatusCode());

        $display = $commandTester->getDisplay();

        foreach (RuleSets::getSetDefinitions() as $name => $set) {
            static::assertStringContainsString($name.') '.$set->getDescription(), $display);
        }
    }
}
