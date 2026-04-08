import { useAuth0 } from '@auth0/auth0-react'

const ROLES_CLAIM = 'https://warehouselogistics/roles'

export type AppRole = 'Employee' | 'Manager' | 'Admin'

export interface Permissions {
  role: AppRole | null
  canReadInventory: boolean
  canReadBOL: boolean
  canCreateBOL: boolean
  canModifyBOL: boolean
  canManageUsers: boolean
}

const PERMISSION_MAP: Record<AppRole, Omit<Permissions, 'role'>> = {
  Employee: {
    canReadInventory: true,
    canReadBOL:       true,
    canCreateBOL:     false,
    canModifyBOL:     false,
    canManageUsers:   false,
  },
  Manager: {
    canReadInventory: true,
    canReadBOL:       true,
    canCreateBOL:     true,
    canModifyBOL:     true,
    canManageUsers:   false,
  },
  Admin: {
    canReadInventory: true,
    canReadBOL:       true,
    canCreateBOL:     true,
    canModifyBOL:     true,
    canManageUsers:   true,
  },
}

const NO_PERMISSIONS: Permissions = {
  role:             null,
  canReadInventory: false,
  canReadBOL:       false,
  canCreateBOL:     false,
  canModifyBOL:     false,
  canManageUsers:   false,
}

export function usePermissions(): Permissions {
  const { user, isAuthenticated } = useAuth0()

  if (!isAuthenticated || !user) return NO_PERMISSIONS

  const roles: string[] = user[ROLES_CLAIM] ?? []
  const role = roles[0] as AppRole | undefined

  if (!role || !(role in PERMISSION_MAP)) return NO_PERMISSIONS

  return { role, ...PERMISSION_MAP[role] }
}
