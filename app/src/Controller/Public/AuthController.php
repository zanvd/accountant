<?php

namespace App\Controller\Public;

use App\Entity\User;
use App\Form\Type\ForgotPasswordType;
use App\Form\Type\LoginType;
use App\Form\Type\RegistrationType;
use App\Form\Type\ResetPasswordType;
use App\Mail\Mailer;
use App\Repository\UserRepository;
use App\Security\EmailVerifier;
use App\Security\ResetPassword;
use Doctrine\Persistence\ManagerRegistry;
use Exception;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\PasswordHasher\Hasher\UserPasswordHasherInterface;
use Symfony\Component\Routing\Annotation\Route;
use SymfonyCasts\Bundle\VerifyEmail\Exception\VerifyEmailExceptionInterface;

#[Route('/', name: 'auth_')]
class AuthController extends AbstractController
{
    private EmailVerifier $emailVerifier;
    private UserRepository $userRepository;

    public function __construct(ManagerRegistry $doctrine, EmailVerifier $emailVerifier)
    {
        $this->emailVerifier = $emailVerifier;
        $this->userRepository = $doctrine->getRepository(User::class);
    }

    #[Route('/forgot-password', methods: ['GET', 'POST'], name: 'forgot_password')]
    public function forgotPassword(Request $request, ResetPassword $resetPassword): Response
    {
        $f = $this->createForm(ForgotPasswordType::class);
        if ($request->isMethod('POST')) {
            $u = $this->userRepository->findOneBy([
                'username' => $request->request->all($f->getName())['username'] ?? ''
            ]);
            if ($u && $resetPassword->initPasswordReset($u)) {
                $this->addFlash('success', 'We\'ve sent a password reset link to your email.');
                
                return $this->redirectToRoute('auth_forgot_password');
            } else {
                $this->addFlash('error', 'Could not send a password reset link to your email. Please, try again.');
            }
        }

        return $this->renderForm('auth/forgot-password.html.twig', ['form' => $f]);
    }

    #[Route('/login', methods: ['GET', 'POST'], name: 'login')]
    public function login(): Response
    {
        $u = new User();
        $f = $this->createForm(LoginType::class, $u);
        
        return $this->renderForm('auth/login.html.twig', ['form' => $f]);
    }

    #[Route('/logout', methods: ['GET'], name: 'logout')]
    public function logout(): void
    {
        throw new Exception('Should not be called.');
    }

    #[Route('/register', methods: ['GET', 'POST'], name: 'register')]
    public function register(UserPasswordHasherInterface $passHasher, Request $request): Response
    {
        $u = new User();
        $f = $this->createForm(RegistrationType::class, $u);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $u->setPassword($passHasher->hashPassword($u, $f->get('password')->getData()));
            $this->userRepository->add($u, true);

            $this->emailVerifier->sendEmailConfirmation($u, 'auth_verify_email');

            return $this->redirectToRoute('dashboard_index');
        }

        return $this->renderForm('auth/register.html.twig', ['form' => $f]);
    }

    #[Route('/reset-password', methods: ['GET', 'POST'], name: 'reset_password')]
    public function resetPassword(
        UserPasswordHasherInterface $passHasher,
        Request $request,
        ResetPassword $resetPassword
    ): Response {
        $t = $request->query->get('token', '');
        $un = $request->query->get('username', '');
        $cun = $resetPassword->getCachedData($t);
        if (!$cun || $un !== $cun) {
            $this->addFlash('error', 'Invalid token provided. Please, request a new link.');

            $this->redirectToRoute('auth_forgot_password');
        }

        $u = $this->userRepository->findOneBy(['username' => $un]);
        if (!$u) {
            $this->addFlash('error', 'Invalid username provided. Please, request a new link.');

            return $this->redirectToRoute('auth_forgot_password');
        }

        $f = $this->createForm(
            ResetPasswordType::class,
            $u,
            [
                'action' => $this->generateUrl(
                    'auth_reset_password',
                    [
                        'token' => $t,
                        'username' => $un,
                    ]
                ),
            ]
        );
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $u->setPassword($passHasher->hashPassword($u, $f->get('password')->getData()));
            $this->userRepository->add($u, true);

            $resetPassword->clearCachedData($t);

            return $this->redirectToRoute('auth_login');
        }

        return $this->renderForm('auth/reset-password.html.twig', ['form' => $f]);
    }

    #[Route('/verify-email', methods: ['GET'], name: 'verify_email')]
    public function verifyUserEmail(Request $request): Response
    {
        $id = $id = $request->get('id');
        if (is_null($id)) return $this->redirectToRoute('auth_register');

        $u = $this->userRepository->find($id);
        if (!$u) return $this->redirectToRoute('auth_register');

        try {
            $this->emailVerifier->handleEmailConfirmation($request, $u);
        } catch (VerifyEmailExceptionInterface $e) {
            $this->addFlash('verify_email_error', $e->getReason());

            return $this->redirectToRoute('auth_register');
        }

        $this->addFlash('success', 'Your email address has been verified.');

        return $this->redirectToRoute('auth_login');
    }
}
