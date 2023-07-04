<?php

namespace App\Exporter;

interface ExporterInterface
{
    public function export(): string;
}
