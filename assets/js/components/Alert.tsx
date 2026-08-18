import type { ReactNode } from 'react';

type AlertProps = {
  title: string;
  children: ReactNode;
  variant?: 'error' | 'info';
};

export function Alert({ title, children, variant = 'info' }: AlertProps) {
  return (
    <div
      className={`perfcheck-alert perfcheck-alert--${variant}`}
      role={variant === 'error' ? 'alert' : 'status'}
    >
      <p className="perfcheck-alert__title">{title}</p>
      <div className="perfcheck-alert__body">{children}</div>
    </div>
  );
}
