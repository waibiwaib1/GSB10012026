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

namespace PhpCsFixer\Console\Command;

use PhpCsFixer\RuleSet\RuleSets;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

/**
 * @author Markus Staab <markus.staab@redaxo.org>
 *
 * @internal
 */
final class ListSetsCommand extends Command
{
    protected static $defaultName = 'list-sets';

    /**
     * {@inheritdoc}
     */
    protected function configure()
    {
        $this
            ->setAliases(['list-rulesets'])
            ->setDescription('List all available built-in rulesets.')
        ;
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        foreach (RuleSets::getSetDefinitionNames() as $name) {
            $output->writeln($name);
        }

        return 0;
    }
}
