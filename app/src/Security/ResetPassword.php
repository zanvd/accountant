<?php

namespace App\Security;

use App\Entity\User;
use App\Mail\Mailer;
use Psr\Cache\CacheItemPoolInterface;
use Symfony\Component\Uid\Uuid;

class ResetPassword
{
    private CacheItemPoolInterface $cache;
    private Mailer $mailer;

    public function __construct(CacheItemPoolInterface $cache, Mailer $mailer)
    {
        $this->cache = $cache;
        $this->mailer = $mailer;
    }

    public function clearCachedData(string $token): void
    {
        $this->cache->deleteItem("auth-pass-reset_$token");
    }

    public function getCachedData(string $token): string
    {
        return $this->cache->getItem("auth-pass-reset_$token")->get() ?? '';
    }

    public function initPasswordReset(User $user): bool
    {
        $t = Uuid::v4();
        $i = $this->cache->getItem("auth-pass-reset_$t");
        $i->set($user->getUsername());
        $i->expiresAfter(3600);
        if ($this->cache->save($i)) {
            $this->mailer->sendPasswordResetEmail($t, $user);

            return true;
        }

        return false;
    }
}
