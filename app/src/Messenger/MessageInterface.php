<?php

namespace App\Messenger;

interface MessageInterface
{
    public function getContent(): int|string;
}
