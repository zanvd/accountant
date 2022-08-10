<?php

namespace App\Form\Type;

use App\Entity\Category;
use App\Entity\TransactionTemplate;
use App\Enum\TransactionType;
use App\Form\Type\Element\ResetIconButtonType;
use App\Form\Type\Element\SubmitIconButtonType;
use Symfony\Bridge\Doctrine\Form\Type\EntityType;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\EnumType;
use Symfony\Component\Form\Extension\Core\Type\FormType;
use Symfony\Component\Form\Extension\Core\Type\IntegerType;
use Symfony\Component\Form\Extension\Core\Type\TextType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolver;

class TransactionTemplateType extends AbstractType
{
    public function buildForm(FormBuilderInterface $builder, array $options): void
    {
        $builder
            ->add('name', TextType::class, [
                'attr' => ['placeholder' => 'Name'],
                'label' => 'Name:',
                'row_attr' => ['class' => 'form-group'],
            ])
            ->add('category', EntityType::class, [
                'choice_label' => 'name',
                'class' => Category::class,
                'label' => 'Category:',
                'row_attr' => ['class' => 'form-group'],
            ])
            ->add('transactionType', EnumType::class, [
                'attr' => ['class' => 'radio-container'],
                'class' => TransactionType::class,
                'expanded' => true,
                'label' => 'Type:',
                'row_attr' => ['class' => 'form-group form-radio'],
            ])
            ->add('position', IntegerType::class, [
                'attr' => [
                    'min' => 0,
                    'placeholder' => 'Position',
                ],
                'label' => 'Position:',
                'row_attr' => ['class' => 'form-group'],
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
        $resolver->setDefaults(['data_class' => TransactionTemplate::class]);
    }
}
