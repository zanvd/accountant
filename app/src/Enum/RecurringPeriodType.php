<?php

namespace App\Enum;

enum RecurringPeriodType: string
{
    case Day = 'day';
    case Week = 'week';
    case Month = 'month';
    case Year = 'year';

    /**
     * Returns a single character symbol compatible with the {@link https://www.php.net/manual/en/dateinterval.construct.php DateInterval's constructor}.
     * @return string
     */
    public function getIntervalSymbol(): string
    {
        return match ($this) {
            RecurringPeriodType::Day => 'D',
            RecurringPeriodType::Week => 'W',
            RecurringPeriodType::Month => 'M',
            RecurringPeriodType::Year => 'Y',
        };
    }
}
