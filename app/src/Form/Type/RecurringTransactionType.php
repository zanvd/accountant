<?php

namespace App\Form\Type;

use App\Entity\Category;
use App\Entity\RecurringTransaction;
use App\Enum\RecurringPeriodType;
use App\Form\Type\Element\ResetIconButtonType;
use App\Form\Type\Element\SubmitIconButtonType;
use Symfony\Bridge\Doctrine\Form\Type\EntityType;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\DateType;
use Symfony\Component\Form\Extension\Core\Type\EnumType;
use Symfony\Component\Form\Extension\Core\Type\FormType;
use Symfony\Component\Form\Extension\Core\Type\MoneyType;
use Symfony\Component\Form\Extension\Core\Type\NumberType;
use Symfony\Component\Form\Extension\Core\Type\TextType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolver;

class RecurringTransactionType extends AbstractType
{
    public function buildForm(FormBuilderInterface $builder, array $options): void
    {
        $builder
            ->add('name', TextType::class, [
                'attr' => ['placeholder' => 'Name'],
                'label' => 'Name:',
            ])
            ->add('category', EntityType::class, [
                'choice_label' => 'name',
                'class' => Category::class,
                'label' => 'Category:',
            ])
            ->add('amount', MoneyType::class, [
                'attr' => [
                    'placeholder' => 'Amount',
                    'step' => '0.01',
                ],
                'currency' => false,
                'help' => 'User . (dot) as a decimal separator.',
                'html5' => true,
                'label' => 'Amount:',
            ])
            ->add('summary', TextType::class, [
                'attr' => ['placeholder' => 'Summary'],
                'empty_data' => '',
                'label' => 'Summary:',
                'required' => false,
            ])
            ->add('periodNum', NumberType::class, [
                'attr' => [
                    'min' => 1,
                ],
                'html5' => true,
                'label' => 'Period num:',
                'scale' => 0,
            ])
            ->add('periodType', EnumType::class, [
                'class' => RecurringPeriodType::class,
                'label' => 'Period type:',
            ])
            ->add('startDate', DateType::class, [
                'format' => 'dd. MM. yyyy',
                'html5' => false,
                'label' => 'Starting on:',
                'widget' => 'single_text',
            ])
            ->add('endDate', DateType::class, [
                'format' => 'dd. MM. yyyy',
                'help' => 'If left empty, it will repeat indefinitely.',
                'html5' => false,
                'label' => 'Ending on:',
                'widget' => 'single_text',
            ])
            ->add(
                $builder->create('actions', FormType::class, [
                    'attr' => ['class' => 'form-actions'],
                    'inherit_data' => true,
                    'label' => false,
                ])
                    ->add('save', SubmitIconButtonType::class, [
                        'attr' => ['class' => 'btn btn-icon'],
                    ])
                    ->add('cancel', ResetIconButtonType::class, [
                        'attr' => ['class' => 'btn btn-icon btn-cancel go-back'],
                    ])
            );
    }
    public function configureOptions(OptionsResolver $resolver): void
    {
        $resolver->setDefaults(['data_class' => RecurringTransaction::class]);
    }
}
