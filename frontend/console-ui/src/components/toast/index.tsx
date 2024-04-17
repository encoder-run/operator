import { alpha, useTheme } from '@mui/material/styles';
import type { FC } from 'react';
import { Toaster } from 'react-hot-toast';

export const Toastr: FC = () => {
  const theme = useTheme();

  return (
    <Toaster
      gutter={24}
      position="top-right"
      reverseOrder={true}
      toastOptions={{
        style: {
          backdropFilter: 'blur(3px)',
          fontWeight: 500,
          padding: theme.spacing(2),
          background: alpha(theme.palette.background.paper, 0.9),
          boxShadow: theme.shadows[21],
          border: 0,
        },
      }}
      containerStyle={{
        top: 24,
        left: 24,
        bottom: 24,
        right: 24,
      }}
    />
  );
};
