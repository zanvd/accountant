<?php

namespace App\Form\Type\Element;

use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\ButtonTypeInterface;
use Symfony\Component\Form\Extension\Core\Type\ButtonType;

class ResetIconButtonType extends AbstractType implements ButtonTypeInterface
{
    public function getParent(): string
    {
        return ButtonType::class;
    }
}
