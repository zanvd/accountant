<?php

namespace App\Form\Type;

use App\Entity\User;
use App\Enum\UserValidationType;
use App\Form\Type\Element\SubmitIconButtonType;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\PasswordType;
use Symfony\Component\Form\Extension\Core\Type\TextType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolver;

class LoginType extends AbstractType
{
    public function buildForm(FormBuilderInterface $builder, array $options): void
    {
        $builder
            ->add('username', TextType::class, [
                'attr' => ['placeholder' => 'Username'],
                'label' => 'Username:',
            ])
            ->add('password', PasswordType::class, [
                'attr' => ['placeholder' => 'Password'],
                'label' => 'Password:',
            ])
            ->add('save', SubmitIconButtonType::class, [
                'attr' => ['class' => 'btn btn-icon'],
                'row_attr' => ['class' => 'form-actions'],
            ]);
    }

    public function configureOptions(OptionsResolver $resolver): void
    {
        $resolver->setDefaults([
            'data_class' => User::class,
            'validation_groups' => [UserValidationType::Login->value],
        ]);
    }
}
