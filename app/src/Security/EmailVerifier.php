<?php

namespace App\Security;

use App\Entity\User;
use App\Mail\Mailer;
use Doctrine\ORM\EntityManagerInterface;
use Symfony\Component\HttpFoundation\Request;
use SymfonyCasts\Bundle\VerifyEmail\Exception\VerifyEmailExceptionInterface;
use SymfonyCasts\Bundle\VerifyEmail\VerifyEmailHelperInterface;

class EmailVerifier
{
    private EntityManagerInterface $entityManager;
    private Mailer $mailer;
    private VerifyEmailHelperInterface $verifyEmailHelper;

    public function __construct(VerifyEmailHelperInterface $helper, Mailer $mailer, EntityManagerInterface $manager)
    {
        $this->entityManager = $manager;
        $this->mailer = $mailer;
        $this->verifyEmailHelper = $helper;
    }

    public function sendEmailConfirmation(User $user, string $verifyRoute): void
    {
        $this->mailer->sendWelcomeEmail(
            $user,
            $this->verifyEmailHelper
                ->generateSignature(
                    $verifyRoute,
                    $user->getId(),
                    $user->getEmail(),
                    ['id' => $user->getId()]
                )
                ->getSignedUrl()
        );
    }

    /**
     * @throws VerifyEmailExceptionInterface
     */
    public function handleEmailConfirmation(Request $request, User $user): void
    {
        $this->verifyEmailHelper->validateEmailConfirmation($request->getUri(), $user->getId(), $user->getEmail());

        $user->setIsVerified(true);

        $this->entityManager->persist($user);
        $this->entityManager->flush();
    }
}
