<?php

namespace App\Mail;

use App\Entity\User;
use Symfony\Component\Mailer\MailerInterface;
use Symfony\Component\Mime\Email;
use Symfony\Component\Routing\Generator\UrlGeneratorInterface;

class Mailer
{
    // Constants represent template IDs in Sendinblue.
    public const MAIL_CONFIRM = 5;
    public const MAIL_RESET = 4;
    public const MAIL_WELCOME = 1;

    private MailerInterface $mailer;
    private UrlGeneratorInterface $router;

    public function __construct(MailerInterface $mailer, UrlGeneratorInterface $router)
    {
        $this->mailer = $mailer;
        $this->router = $router;
    }

    public function sendEmail(string $subject, int $templateId, string $to, array $params = []): void
    {
        $email = new Email();
        $email->to($to)
            ->subject($subject)
            ->text('dummy')
            ->getHeaders()
            ->addTextHeader('templateId', $templateId);
        if ($params) $email->getHeaders()->addParameterizedHeader('params', 'params', $params);
        $this->mailer->send($email);
    }

    public function sendPasswordResetEmail(string $token, User $user): void
    {
        $this->sendEmail(
            'Reset password',
            static::MAIL_RESET,
            $user->getEmail(),
            [
                'url' => $this->router->generate(
                    'auth_reset_password',
                    [
                        'token' => $token,
                        'username' => $user->getUsername(),
                    ],
                    UrlGeneratorInterface::ABSOLUTE_URL
                ),
                'username' => $user->getUsername(),
            ]
        );
    }

    public function sendWelcomeEmail(User $user, string $signedUrl): void
    {
        $this->sendEmail(
            'Welcome to Accountant',
            static::MAIL_WELCOME,
            $user->getEmail(),
            [
                'url' => $signedUrl,
                'username' => $user->getUsername(),
            ]
        );
    }
}
