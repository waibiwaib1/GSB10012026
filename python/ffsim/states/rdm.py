# (C) Copyright IBM 2023.
#
# This code is licensed under the Apache License, Version 2.0. You may
# obtain a copy of this license in the LICENSE.txt file in the root directory
# of this source tree or at http://www.apache.org/licenses/LICENSE-2.0.
#
# Any modifications or derivative works of this code must retain this
# copyright notice, and modified files need to carry a notice indicating
# that they have been altered from the originals.

"""Functions to compute reduced density matrices."""

from __future__ import annotations

import numpy as np
from pyscf.fci.direct_spin1 import (
    make_rdm1,
    make_rdm1s,
    make_rdm12,
    make_rdm12s,
    trans_rdm1,
    trans_rdm1s,
    trans_rdm12,
    trans_rdm12s,
)

from ffsim.cistring import gen_linkstr_index


def rdm(
    vec: np.ndarray,
    norb: int,
    nelec: tuple[int, int],
    rank: int = 1,
    spin_summed: bool = True,
    reordered: bool = True,
    return_lower_ranks: bool = True,
) -> np.ndarray | tuple:
    """Return the reduced density matrix (RDM) or matrices of a state vector.

    The rank 1 RDM is defined as follows:

    .. code::

        rdm1[p, q] = ⟨p+ q⟩

    The definition of higher-rank RDMs depends on the ``reordered`` argument, which
    defaults to True.

    **reordered = True**

    The reordered RDMs are defined as follows:

    .. code::

        rdm2[p, q, r, s] = ⟨p+ r+ s q⟩
        rdm3[p, q, r, s, t, u] = ⟨p+ r+ t+ u s q⟩
        rdm4[p, q, r, s, t, u, v, w] = ⟨p+ r+ t+ v+ w u s q⟩

    **reordered = False**

    If reordered is set to False, the RDMs are defined as follows:

    .. code::

        rdm2[p, q, r, s] = ⟨p+ q r+ s⟩
        rdm3[p, q, r, s, t, u] = ⟨p+ q r+ s t+ u⟩
        rdm4[p, q, r, s, t, u, v, w] = ⟨p+ q r+ s t+ u v+ w⟩

    Note:
        Currently, only ranks 1 and 2 are supported.

    Args:
        vec: The state vector whose reduced density matrix is desired.
        norb: The number of spatial orbitals.
        nelec: The number of alpha and beta electrons.
        rank: The rank of the reduced density matrix.
        spin_summed: Whether to sum over the spin index.
        reordered: Whether to reorder the indices of the reduced density matrix.
        return_lower_ranks: Whether to return lower rank RDMs in addition to the
            specified rank. If True, then this function returns all RDMs up to and
            including the specified rank, in increasing order of rank. For example,
            if `rank=2` then a tuple `(rdm1, rdm2)` is returned.

    Returns:
        The reduced density matrix or matrices. If `return_lower_ranks` is False,
        then a single RDM is returned. If `return_lower_ranks` is True, then
        a `rank`-length tuple of RDMs is returned, containing the RDMs up to
        the specified rank in increasing order of rank.

        If `spin_summed` is True, then each RDM is a single array with only
        spatial orbital indices. If `spin_summed` is False, then the spin
        alpha and beta sectors are not combined into a single array. Instead,
        the rank 1 RDM is returned as a pair `(rdm_a, rdm_b)`, and the rank 2
        RDM is returned as a tuple `(rdm_aa, rdm_ab, rdm_ba, rdm_bb)`, in the
        style of PySCF's `make_rdm1s` and `make_rdm12s` functions.
    """
    n_alpha, n_beta = nelec
    link_index_a = gen_linkstr_index(range(norb), n_alpha)
    link_index_b = gen_linkstr_index(range(norb), n_beta)
    link_index = (link_index_a, link_index_b)
    if rank == 1:
        if spin_summed:
            return _rdm1_spin_summed(vec, norb, nelec, link_index)
        else:
            return _rdm1(vec, norb, nelec, link_index)
    if rank == 2:
        if spin_summed:
            return _rdm2_spin_summed(
                vec, norb, nelec, reordered, link_index, return_lower_ranks
            )
        else:
            return _rdm2(vec, norb, nelec, reordered, link_index, return_lower_ranks)
    raise NotImplementedError(
        f"Computing the rank {rank} reduced density matrix is currently not supported."
    )


def _rdm1_spin_summed(
    vec: np.ndarray,
    norb: int,
    nelec: tuple[int, int],
    link_index: tuple[np.ndarray, np.ndarray] | None,
) -> np.ndarray:
    rdm1_real = make_rdm1(vec.real, norb, nelec, link_index=link_index)
    rdm1_imag = make_rdm1(vec.imag, norb, nelec, link_index=link_index)
    trans_rdm1_real_imag = trans_rdm1(
        vec.real, vec.imag, norb, nelec, link_index=link_index
    )
    trans_rdm1_imag_real = trans_rdm1(
        vec.imag, vec.real, norb, nelec, link_index=link_index
    )
    return _assemble_rdm1_spin_summed(
        rdm1_real, rdm1_imag, trans_rdm1_real_imag, trans_rdm1_imag_real
    )


def _rdm1(
    vec: np.ndarray,
    norb: int,
    nelec: tuple[int, int],
    link_index: tuple[np.ndarray, np.ndarray] | None,
) -> np.ndarray:
    rdms1_real = make_rdm1s(vec.real, norb, nelec, link_index=link_index)
    rdms1_imag = make_rdm1s(vec.imag, norb, nelec, link_index=link_index)
    trans_rdms1_real_imag = trans_rdm1s(
        vec.real, vec.imag, norb, nelec, link_index=link_index
    )
    trans_rdms1_imag_real = trans_rdm1s(
        vec.imag, vec.real, norb, nelec, link_index=link_index
    )
    return _assemble_rdm1(
        rdms1_real, rdms1_imag, trans_rdms1_real_imag, trans_rdms1_imag_real
    )


def _rdm2_spin_summed(
    vec: np.ndarray,
    norb: int,
    nelec: tuple[int, int],
    reordered: bool,
    link_index: tuple[np.ndarray, np.ndarray] | None,
    return_lower_ranks: bool,
) -> np.ndarray | tuple[np.ndarray, np.ndarray]:
    rdm1_real, rdm2_real = make_rdm12(
        vec.real, norb, nelec, reorder=reordered, link_index=link_index
    )
    rdm1_imag, rdm2_imag = make_rdm12(
        vec.imag, norb, nelec, reorder=reordered, link_index=link_index
    )
    trans_rdm1_real_imag, trans_rdm2_real_imag = trans_rdm12(
        vec.real, vec.imag, norb, nelec, reorder=reordered, link_index=link_index
    )
    trans_rdm1_imag_real, trans_rdm2_imag_real = trans_rdm12(
        vec.imag, vec.real, norb, nelec, reorder=reordered, link_index=link_index
    )
    rdm2 = _assemble_rdm2_spin_summed(
        rdm2_real, rdm2_imag, trans_rdm2_real_imag, trans_rdm2_imag_real
    )
    if not return_lower_ranks:
        return rdm2
    rdm1 = _assemble_rdm1_spin_summed(
        rdm1_real, rdm1_imag, trans_rdm1_real_imag, trans_rdm1_imag_real
    )
    return rdm1, rdm2


def _rdm2(
    vec: np.ndarray,
    norb: int,
    nelec: tuple[int, int],
    reordered: bool,
    link_index: tuple[np.ndarray, np.ndarray] | None,
    return_lower_ranks: bool,
) -> np.ndarray | tuple[np.ndarray, np.ndarray]:
    rdms1_real, rdms2_real = make_rdm12s(
        vec.real, norb, nelec, reorder=reordered, link_index=link_index
    )
    rdms1_imag, rdms2_imag = make_rdm12s(
        vec.imag, norb, nelec, reorder=reordered, link_index=link_index
    )
    trans_rdms1_real_imag, trans_rdms2_real_imag = trans_rdm12s(
        vec.real, vec.imag, norb, nelec, reorder=reordered, link_index=link_index
    )
    trans_rdms1_imag_real, trans_rdms2_imag_real = trans_rdm12s(
        vec.imag, vec.real, norb, nelec, reorder=reordered, link_index=link_index
    )
    rdm2 = _assemble_rdm2(
        rdms2_real, rdms2_imag, trans_rdms2_real_imag, trans_rdms2_imag_real
    )
    if not return_lower_ranks:
        return rdm2
    rdm1 = _assemble_rdm1(
        rdms1_real, rdms1_imag, trans_rdms1_real_imag, trans_rdms1_imag_real
    )
    return rdm1, rdm2


def _assemble_rdm1_spin_summed(
    rdm1_real: np.ndarray,
    rdm1_imag: np.ndarray,
    trans_rdm1_real_imag: np.ndarray,
    trans_rdm1_imag_real: np.ndarray,
) -> np.ndarray:
    # use minus sign for 1j because the rdm1 convention from pyscf is transposed
    return rdm1_real + rdm1_imag - 1j * (trans_rdm1_real_imag - trans_rdm1_imag_real)


def _assemble_rdm1(
    rdms1_real: tuple[np.ndarray, np.ndarray],
    rdms1_imag: tuple[np.ndarray, np.ndarray],
    trans_rdms1_real_imag: tuple[np.ndarray, np.ndarray],
    trans_rdms1_imag_real: tuple[np.ndarray, np.ndarray],
) -> tuple[np.ndarray, np.ndarray]:
    rdm1_a, rdm1_b = (rdm.astype(complex) for rdm in rdms1_real)
    rdm1_a += rdms1_imag[0]
    rdm1_b += rdms1_imag[1]
    # use minus sign for real-imag and plus sign for imag-real because
    # the rdm1 convention in pyscf is transposed
    rdm1_a -= 1j * trans_rdms1_real_imag[0]
    rdm1_b -= 1j * trans_rdms1_real_imag[1]
    rdm1_a += 1j * trans_rdms1_imag_real[0]
    rdm1_b += 1j * trans_rdms1_imag_real[1]
    return rdm1_a, rdm1_b


def _assemble_rdm2_spin_summed(
    rdm2_real: np.ndarray,
    rdm2_imag: np.ndarray,
    trans_rdm2_real_imag: np.ndarray,
    trans_rdm2_imag_real: np.ndarray,
) -> np.ndarray:
    return rdm2_real + rdm2_imag + 1j * (trans_rdm2_real_imag - trans_rdm2_imag_real)


def _assemble_rdm2(
    rdms2_real: tuple[np.ndarray, np.ndarray, np.ndarray],
    rdms2_imag: tuple[np.ndarray, np.ndarray, np.ndarray],
    trans_rdms2_real_imag: tuple[np.ndarray, np.ndarray, np.ndarray, np.ndarray],
    trans_rdms2_imag_real: tuple[np.ndarray, np.ndarray, np.ndarray, np.ndarray],
) -> tuple[np.ndarray, np.ndarray, np.ndarray, np.ndarray]:
    rdm_aa, rdm_ab, rdm_bb = (rdm.astype(complex) for rdm in rdms2_real)
    rdm_aa += rdms2_imag[0]
    rdm_ab += rdms2_imag[1]
    rdm_bb += rdms2_imag[2]
    rdm_aa += 1j * trans_rdms2_real_imag[0]
    rdm_ab += 1j * trans_rdms2_real_imag[1]
    rdm_bb += 1j * trans_rdms2_real_imag[3]
    rdm_aa -= 1j * trans_rdms2_imag_real[0]
    rdm_ab -= 1j * trans_rdms2_imag_real[1]
    rdm_bb -= 1j * trans_rdms2_imag_real[3]
    # rdms2 is currently (rdm_aa, rdm_ab, rdm_bb)
    # the rdm_ba block is obtained by transposing the rdm_ab block
    rdm_ba = rdm_ab.transpose(2, 3, 0, 1)
    return rdm_aa, rdm_ab, rdm_ba, rdm_bb
