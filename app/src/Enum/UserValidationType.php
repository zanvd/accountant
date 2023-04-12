<?php

namespace App\Enum;

enum UserValidationType: string
{
    case Login = 'login';
    case ForgotPassword = 'forgot_password';
    case Register = 'register';
    case Reset = 'reset';
}
